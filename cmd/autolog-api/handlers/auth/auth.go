package auth

import (
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type AuthHandler struct {
	// foundationals/platform
	calendarService calendar.ServiceIface

	publicKey []byte

	// services
	userService user.ServiceIface
}

type AuthHandlerConfig struct {
	// foundationals/platform
	CalendarService calendar.ServiceIface
	RandomGenerator random.ServiceIface
	Logger          *logger.Logger

	PublicKey []byte

	// services
	UserService user.ServiceIface
}

func NewAuthHandler(config AuthHandlerConfig) *AuthHandler {
	return &AuthHandler{
		calendarService: config.CalendarService,
		publicKey:       config.PublicKey,

		userService: config.UserService,
	}
}
