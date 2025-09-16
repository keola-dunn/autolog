package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type TokenVerifier struct {
	keyFunc keyfunc.Keyfunc
}

type TokenVerifierConfig struct {
	JWKSUrl string
}

func NewTokenVerifier(ctx context.Context, config TokenVerifierConfig) (*TokenVerifier, error) {
	jwksFunc, err := keyfunc.NewDefaultCtx(ctx, []string{
		config.JWKSUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create jwksFunc: %w", err)
	}

	verifier := TokenVerifier{
		keyFunc: jwksFunc,
	}

	return &verifier, nil
}

func (v *TokenVerifier) VerifyToken(tokenString string) (bool, AutologAPIJWTClaims, error) {
	var claims AutologAPIJWTClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, v.keyFunc.Keyfunc)
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
