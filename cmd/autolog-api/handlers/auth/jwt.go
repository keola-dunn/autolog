package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/keola-dunn/autolog/cmd/autolog-api/utils"
	"github.com/keola-dunn/autolog/internal/httputil"

	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthHandler) createJWT(userId string) (string, error) {
	now := h.calendarService.NowUTC()

	tokenId, err := h.randomGenerator.RandomUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate random uuid token id: %w", err)
	}

	claims := jwt.RegisteredClaims{
		Issuer:    h.jwtIssuer,
		Subject:   userId,
		Audience:  jwt.ClaimStrings{}, // app specific keys indicating what the JWT is intended to be used by
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(h.jwtExpiryLengthMinutes))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        tokenId,
	}

	myClaims := utils.AutologAPIJWTClaims{
		RegisteredClaims: claims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)

	jwtToken, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return jwtToken, nil
}

// VerifyToken makes sure the token is valid. Returns boolean indicating if the token
// is valid, the user id associated with the token, and an error
func (a *AuthHandler) VerifyToken(tokenString string) (bool, utils.AutologAPIJWTClaims, error) {
	var claims utils.AutologAPIJWTClaims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})
	if err != nil {
		return false, claims, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if !token.Valid {
		return false, claims, nil
	}

	return true, claims, nil
}

// RequireAuthentication is a middleware that requires the request to be authenticated
func (a *AuthHandler) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		token := splitToken[1]

		valid, _, err := a.VerifyToken(token)
		if err != nil {
			a.logger.Error("failed to verify token", err)
			httputil.RespondWithError(w, http.StatusInternalServerError, "")
			return
		}
		if !valid {
			// log failed auth attempts
			a.logger.Warn("invalid token provided",
				"token", token,
				"referer", r.Header.Get("referer"),
				"user-agent", r.Header.Get("user-agent"),
				"x-forwarded-for", r.Header.Get("X-Forwarded-For"))

			httputil.RespondWithError(w, http.StatusUnauthorized, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}
