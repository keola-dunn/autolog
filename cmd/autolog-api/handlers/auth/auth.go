package auth

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type AuthHandler struct {
	// foundationals/platform
	calendarService calendar.ServiceIface

	publicKeyData []byte
	publicKey     *rsa.PublicKey

	// services
	userService user.ServiceIface
}

type AuthHandlerConfig struct {
	// foundationals/platform
	CalendarService calendar.ServiceIface
	RandomGenerator random.ServiceIface
	Logger          *logger.Logger

	PublicKeyData []byte

	// services
	UserService user.ServiceIface
}

func NewAuthHandler(config AuthHandlerConfig) (*AuthHandler, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(config.PublicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return &AuthHandler{
		calendarService: config.CalendarService,

		userService: config.UserService,

		publicKeyData: config.PublicKeyData,
		publicKey:     pubKey,
	}, nil
}
