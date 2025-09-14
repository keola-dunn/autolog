package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokenFromAuthHeader(authHeader string) string {
	if strings.TrimSpace(authHeader) == "" {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 || !strings.Contains(authHeader, "Bearer") {
		return ""
	}

	return splitToken[1]
}

type AutologAPIJWTClaims struct {
	jwt.RegisteredClaims
}

func (a *AutologAPIJWTClaims) GetUserId() string {
	return a.Subject
}

type VerifyTokenInput struct {
	TokenString string

	PublicKey []byte

	WellKnownJWKSUrl string
}

// VerifyToken makes sure the token is valid. Returns boolean indicating if the token
// is valid, the user id associated with the token, and an error
func VerifyToken(tokenString string, publicKey *rsa.PublicKey) (bool, AutologAPIJWTClaims, error) {
	var claims AutologAPIJWTClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, claims, jwt.ErrTokenExpired
		}

		return false, claims, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if !token.Valid {
		return false, claims, nil
	}

	return true, claims, nil
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

	// TokenSecret is the private key used to sign the token. This is not a public value.
	PrivateKey []byte
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

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(input.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, myClaims)

	jwtToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return jwtToken, nil
}
