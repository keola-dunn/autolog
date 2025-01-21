package auth

import "github.com/keola-dunn/autolog/internal/platform/postgres"

type ServiceConfig struct {
	DB postgres.ConnectionPool
}

type Service struct {
	db postgres.ConnectionPool
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db: cfg.DB,
	}
}
