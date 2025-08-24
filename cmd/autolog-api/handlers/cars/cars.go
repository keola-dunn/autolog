package cars

import (
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/car"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type CarsHandler struct {
	// foundationals/platform
	calendarService calendar.ServiceIface
	randomGenerator random.ServiceIface
	logger          *logger.Logger

	// services
	userService user.ServiceIface
	carsHandler car.ServiceIface

	jwtSecret string
}

type CarsHandlerConfig struct {
	// foundationals/platform
	CalendarService calendar.ServiceIface
	RandomGenerator random.ServiceIface
	Logger          *logger.Logger

	// services
	UserService user.ServiceIface
	CarsHandler car.ServiceIface

	JWTSecret string
}

func NewCarsHandler(config CarsHandlerConfig) *CarsHandler {
	return &CarsHandler{
		calendarService: config.CalendarService,
		randomGenerator: config.RandomGenerator,
		logger:          config.Logger,

		userService: config.UserService,
		carsHandler: config.CarsHandler,
		jwtSecret:   config.JWTSecret,
	}
}
