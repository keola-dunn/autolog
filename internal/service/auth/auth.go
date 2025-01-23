package auth

import (
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
)

type ServiceConfig struct {
	DB              postgres.ConnectionPool
	RandomGenerator random.ServiceIface
}

type Service struct {
	db              postgres.ConnectionPool
	randomGenerator random.ServiceIface
	saltLength      int64
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.RandomGenerator == nil {
		cfg.RandomGenerator = random.NewService()
	}

	return &Service{
		db:              cfg.DB,
		randomGenerator: cfg.RandomGenerator,
		saltLength:      8,
	}
}
