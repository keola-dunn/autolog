package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AutologAPIJWTClaims struct {
	jwt.RegisteredClaims
}

// VerifyToken makes sure the token is valid. Returns boolean indicating if the token
// is valid, the user id associated with the token, and an error
func VerifyToken(tokenString, jwtSecret string) (bool, AutologAPIJWTClaims, error) {
	var claims AutologAPIJWTClaims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return false, claims, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if !token.Valid {
		return false, claims, nil
	}

	return true, claims, nil
}

func GetUserId(claims AutologAPIJWTClaims) string {
	return claims.Subject
}

type CreateJWTInput struct {
	// Issuer is the service that created and issued the token
	Issuer string

	// UserId is the ID of the user whose token this is. This will be the token subject.
	UserId string

	// IssuedAt is when the token was created
	IssuedAt time.Time

	// ExpiresAt is when the token is expired
	ExpiresAt time.Time

	// NotBefore is when the token can be used
	NotBefore time.Time

	// Id is the ID of the token
	Id string

	// TokenSecret is the secret value used to sign the token. This is not a public value.
	TokenSecret string
}

func CreateJWT(input CreateJWTInput) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    input.Issuer,
		Subject:   input.UserId,
		Audience:  jwt.ClaimStrings{}, // app specific keys indicating what the JWT is intended to be used by
		ExpiresAt: jwt.NewNumericDate(input.ExpiresAt),
		NotBefore: jwt.NewNumericDate(input.NotBefore),
		IssuedAt:  jwt.NewNumericDate(input.IssuedAt),
		ID:        input.Id,
	}

	myClaims := AutologAPIJWTClaims{
		RegisteredClaims: claims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)

	jwtToken, err := token.SignedString([]byte(input.TokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return jwtToken, nil
}
