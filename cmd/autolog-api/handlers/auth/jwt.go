package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AutologAPIJWTClaims struct {
	jwt.RegisteredClaims
}

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

	myClaims := AutologAPIJWTClaims{
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
// is valid, the user id associated with the token,
func (a *AuthHandler) VerifyToken(tokenString string) (bool, AutologAPIJWTClaims, error) {
	var claims AutologAPIJWTClaims
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

// func (a *AuthHandler) JWTAuthention(next http.Handler) http.Handler {

// }
