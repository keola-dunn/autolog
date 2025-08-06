package shop

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	CreateShop(ctx context.Context, shop Shop, creatorUserId string) (string, error)
	SearchForShop(ctx context.Context, term string, limit int64) ([]Shop, error)
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

type Shop struct {
	id       string
	Name     string
	Address1 string
	Address2 string
	City     string
	State    string
	Zip      string
	Phone    string

	createdBy string
	createdAt time.Time
	updatedAt time.Time
}

func (s *Shop) valid() bool {
	// TODO: build into this at a later date
	return strings.TrimSpace(s.Name) != ""
}

func (s *Shop) Id() string {
	return s.id
}

func (s *Shop) CreatedBy() string {
	return s.createdBy
}

func (s *Shop) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Shop) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Service) CreateShop(ctx context.Context, shop Shop, creatorUserId string) (string, error) {
	if s.db == nil {
		return "", ErrMissingRequiredConfiguration
	}

	if !shop.valid() {
		return "", ErrInvalidArg
	}

	query := `
	INSERT INTO shops (
    	name,
    	address1,
    	address2,
    	city,
    	state,
    	zip,
    	phone,
    	created_by
	) VALUES 
	 ($1, $2, $3, $4, $5, $6, $7, $8)
	 RETURNING id`

	row := s.db.QueryRow(ctx, query, shop.Name, shop.Address1,
		shop.Address2, shop.City, shop.State, shop.Zip,
		shop.Phone, creatorUserId)

	var shopId string
	if err := row.Scan(&shopId); err != nil {
		return "", fmt.Errorf("failed to insert shop: %w", err)
	}

	return shopId, nil
}

func (s *Service) SearchForShop(ctx context.Context, term string, limit int64) ([]Shop, error) {
	if s.db == nil {
		return nil, ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(term) == "" || len(term) <= 3 {
		return nil, ErrInvalidArg
	}

	query := `
	SELECT 
		s.id,
    	s.name,
    	s.address1,
    	s.address2,
    	s.city,
    	s.state,
    	s.zip,
    	s.phone,
    	s.created_by,
    	s.created_at,
    	s.updated_at
	FROM shops s
	WHERE 
		LOWER(CONCAT(s.name, ' - ', s.city)) LIKE LOWER(%1) OR 
		LOWER(CONCAT(s.name, ' - ', s.address1, ' ', s.address2, ', ', s.city, ', ', s.state, ' ', s.zip)) LIKE LOWER(%1)
	ORDER BY CONCAT(s.name, ' - ', s.city)
	LIMIT $2`

	rows, err := s.db.Query(ctx, query, strings.TrimSpace(strings.ToLower(fmt.Sprintf("%%%s%%", term))), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query for shops: %w", err)
	}
	defer rows.Close()

	var shops = make([]Shop, 0, limit)
	for rows.Next() {
		var shop Shop

		if err := rows.Scan(&shop); err != nil {
			return nil, fmt.Errorf("failed to scan shop: %w", err)
		}

		shops = append(shops, shop)
	}

	return shops, nil
}
