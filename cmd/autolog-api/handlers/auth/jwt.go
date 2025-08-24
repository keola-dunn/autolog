package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/keola-dunn/autolog/internal/httputil"

	autologjwt "github.com/keola-dunn/autolog/cmd/autolog-api/jwt"
)

func (h *AuthHandler) createJWT(userId string) (string, error) {
	now := h.calendarService.NowUTC()

	tokenId, err := h.randomGenerator.RandomUUID()
	if err != nil {
		return "", fmt.Errorf("failed to create random token id: %w", err)
	}

	jwtToken, err := autologjwt.CreateJWT(autologjwt.CreateJWTInput{
		Issuer:      h.jwtIssuer,
		UserId:      userId,
		IssuedAt:    now,
		ExpiresAt:   now.Add(time.Duration(h.jwtExpiryLengthMinutes) * time.Minute),
		NotBefore:   now,
		Id:          tokenId,
		TokenSecret: h.jwtSecret,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create jwt: %w", err)
	}

	return jwtToken, nil
}

// RequireAuthentication is a middleware that requires the request to be authenticated
func (a *AuthHandler) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		token := splitToken[1]

		valid, _, err := autologjwt.VerifyToken(token, a.jwtSecret)
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
