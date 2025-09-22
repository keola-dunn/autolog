package image

import (
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
)

type Service struct {
	imagePrefix string

	db postgres.ConnectionPool

	randomGenerator random.ServiceIface
}

type ServiceConfig struct {
	ImagePrefix string

	// DB is the Database used for the auth service
	DB postgres.ConnectionPool

	RandomGenerator random.ServiceIface
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		imagePrefix:     cfg.ImagePrefix,
		db:              cfg.DB,
		randomGenerator: cfg.RandomGenerator,
	}
}
