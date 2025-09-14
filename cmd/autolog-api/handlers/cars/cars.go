package cars

import (
	"github.com/keola-dunn/autolog/internal/calendar"
	"github.com/keola-dunn/autolog/internal/logger"
	nhtsavpic "github.com/keola-dunn/autolog/internal/nhtsa"
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
	carService  car.ServiceIface

	nhtsaClient nhtsavpic.ClientIface

	jwtPublicKey []byte
}

type CarsHandlerConfig struct {
	// foundationals/platform
	CalendarService calendar.ServiceIface
	RandomGenerator random.ServiceIface
	Logger          *logger.Logger

	// services
	UserService user.ServiceIface
	CarService  car.ServiceIface

	NHTSAClient nhtsavpic.ClientIface

	JWTPublicKey []byte
}

func NewCarsHandler(config CarsHandlerConfig) *CarsHandler {
	return &CarsHandler{
		calendarService: config.CalendarService,
		randomGenerator: config.RandomGenerator,
		logger:          config.Logger,

		userService: config.UserService,
		carService:  config.CarService,

		nhtsaClient: config.NHTSAClient,

		jwtPublicKey: config.JWTPublicKey,
	}
}
