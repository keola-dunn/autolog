package auth

import (
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
)

type ServiceConfig struct {
	// DB is the Database used for the auth service
	DB postgres.ConnectionPool

	// RandomGenerator is used to generate random values within the auth service.
	RandomGenerator random.ServiceIface

	// SaltLength sets the length of the password salts generated for the auth service.
	// Defaults to 8.
	SaltLength int64
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

	if cfg.SaltLength <= 0 {
		cfg.SaltLength = 8
	}

	return &Service{
		db:              cfg.DB,
		randomGenerator: cfg.RandomGenerator,
		saltLength:      cfg.SaltLength,
	}
}
