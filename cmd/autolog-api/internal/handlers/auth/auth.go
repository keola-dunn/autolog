package auth

import (
	autologjwt "github.com/keola-dunn/autolog/internal/jwt"
)

type AuthHandler struct {
	jwtVerifier *autologjwt.TokenVerifier
}

type AuthHandlerConfig struct {
	// foundationals/platform

	TokenVerifier *autologjwt.TokenVerifier
}

func NewAuthHandler(config AuthHandlerConfig) (*AuthHandler, error) {
	return &AuthHandler{
		jwtVerifier: config.TokenVerifier,
	}, nil
}
