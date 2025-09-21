package auth

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/keola-dunn/autolog/internal/calendar"
	autologjwt "github.com/keola-dunn/autolog/internal/jwt"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type AuthHandler struct {
	autologjwt.AuthHandler

	// configs
	jwtIssuer              string
	jwtExpiryLengthMinutes int64

	// foundationals/platform
	calendarService calendar.ServiceIface
	randomGenerator random.ServiceIface
	logger          *logger.Logger

	// services
	userService user.ServiceIface

	jwtPublicKeyData []byte
	jwtPublicKey     *rsa.PublicKey

	jwtPrivateKeyData []byte
	jwtPrivateKey     *rsa.PrivateKey
}

type AuthHandlerConfig struct {
	JWTIssuer              string
	JWTExpiryLengthMinutes int64

	// foundationals/platform
	CalendarService calendar.ServiceIface
	RandomGenerator random.ServiceIface
	Logger          *logger.Logger

	// services
	UserService user.ServiceIface

	JWTPublicKeyData  []byte
	JWTPrivateKeyData []byte
}

func NewAuthHandler(config AuthHandlerConfig) (*AuthHandler, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(config.JWTPublicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(config.JWTPrivateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &AuthHandler{
		jwtIssuer:              config.JWTIssuer,
		jwtExpiryLengthMinutes: config.JWTExpiryLengthMinutes,

		calendarService: config.CalendarService,
		randomGenerator: config.RandomGenerator,
		logger:          config.Logger,

		userService: config.UserService,

		jwtPublicKeyData:  config.JWTPublicKeyData,
		jwtPublicKey:      pubKey,
		jwtPrivateKeyData: config.JWTPrivateKeyData,
		jwtPrivateKey:     privKey,
	}, nil
}
