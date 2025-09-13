package auth

import (
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type AuthHandler struct {
	// configs
	jwtSecret              string
	jwtIssuer              string
	jwtExpiryLengthMinutes int64

	// foundationals/platform
	calendarService calendar.ServiceIface
	randomGenerator random.ServiceIface
	logger          *logger.Logger

	// services
	userService user.ServiceIface

	jwtPublicKey  []byte
	jwtPrivateKey []byte
}

type AuthHandlerConfig struct {
	// JWTSecret              string
	JWTIssuer              string
	JWTExpiryLengthMinutes int64

	// foundationals/platform
	CalendarService calendar.ServiceIface
	RandomGenerator random.ServiceIface
	Logger          *logger.Logger

	// services
	UserService user.ServiceIface

	JWTPublicKey  []byte
	JWTPrivateKey []byte
}

func NewAuthHandler(config AuthHandlerConfig) *AuthHandler {
	return &AuthHandler{
		//jwtSecret:              config.JWTSecret,
		jwtIssuer:              config.JWTIssuer,
		jwtExpiryLengthMinutes: config.JWTExpiryLengthMinutes,

		calendarService: config.CalendarService,
		randomGenerator: config.RandomGenerator,
		logger:          config.Logger,

		userService: config.UserService,

		jwtPublicKey:  config.JWTPublicKey,
		jwtPrivateKey: config.JWTPrivateKey,
	}
}
