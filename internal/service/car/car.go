package car

import (
	"errors"
	"time"

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
}

type ServiceIface interface {
}

type Service struct {
	db              postgres.ConnectionPool
	randomGenerator random.ServiceIface
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.RandomGenerator == nil {
		cfg.RandomGenerator = random.NewService()
	}

	return &Service{
		db:              cfg.DB,
		randomGenerator: cfg.RandomGenerator,
	}
}

type Car struct {
	Id    string
	Make  string
	Model string
	Trim  string
	Year  int64
	VIN   string

	CreatedAt time.Time
	UpdatedAt time.Time
}
