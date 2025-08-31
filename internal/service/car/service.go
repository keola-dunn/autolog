package car

import (
	"context"
	"errors"

	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
)

var (
	ErrMissingRequiredConfiguration = errors.New("auth service is missing required configurations to perform this operation")

	ErrInvalidArg = errors.New("one or more of the provided arguments are invalid")
)

type ServiceConfig struct {
	// DB is the Database used for the auth service
	DB postgres.ConnectionPool

	// RandomGenerator is used to generate random values within the auth service.
	RandomGenerator random.ServiceIface

	// PublicIdLength is the length of a public id assigned to a car. Defaults to 6.
	PublicIdLength int64
}

type ServiceIface interface {
	CreateServiceLog(ctx context.Context, serviceLog ServiceLog, userId, carId string) (string, error)
	CreateCar(ctx context.Context, userId string, car Car, nhtsaData NHTSAVPICData) error
	GetCar(ctx context.Context, input GetCarInput) (GetCarOutput, error)
}

type Service struct {
	db              postgres.ConnectionPool
	randomGenerator random.ServiceIface

	publicIdLength int64
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.RandomGenerator == nil {
		cfg.RandomGenerator = random.NewService()
	}

	if cfg.PublicIdLength < 6 {
		cfg.PublicIdLength = 6
	}

	return &Service{
		db:              cfg.DB,
		randomGenerator: cfg.RandomGenerator,
		publicIdLength:  cfg.PublicIdLength,
	}
}
