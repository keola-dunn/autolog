package cars

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
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

	jwtPublicKeyData []byte
	publicKey        *rsa.PublicKey
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

	JWTPublicKeyData []byte
}

func NewCarsHandler(config CarsHandlerConfig) (*CarsHandler, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(config.JWTPublicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &CarsHandler{
		calendarService: config.CalendarService,
		randomGenerator: config.RandomGenerator,
		logger:          config.Logger,

		userService: config.UserService,
		carService:  config.CarService,

		nhtsaClient: config.NHTSAClient,

		jwtPublicKeyData: config.JWTPublicKeyData,
		publicKey:        pubKey,
	}, nil
}
