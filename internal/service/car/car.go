package car

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
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
	CreateServiceLog(ctx context.Context, serviceLog ServiceLog, userId, carId string) (string, error)
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
	id    string
	Make  string
	Model string
	Trim  string
	Year  int64
	VIN   string

	createdAt time.Time
	updatedAt time.Time
}

func (c *Car) Id() string {
	return c.id
}

func (s *Service) CreateCar(ctx context.Context, userId string, car Car) error {
	if s.db == nil {
		return ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(userId) == "" || !car.valid() {
		return ErrInvalidArg
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	carId, err := createCarRecord(ctx, tx, car)
	if err != nil {
		return fmt.Errorf("failed to create car redord: %w", err)
	}

	if _, err := createUserCarRecord(ctx, tx, userId, carId); err != nil {
		return fmt.Errorf("failed to create user car record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *Car) valid() bool {
	// first car created was in 1885
	return !(strings.TrimSpace(c.Make) == "" ||
		strings.TrimSpace(c.Model) == "" ||
		c.Year < 1885 ||
		strings.TrimSpace(c.VIN) == "")

}

func createCarRecord(ctx context.Context, tx pgx.Tx, car Car) (string, error) {
	query := `
	INSERT INTO cars (make, model, trim, year, vin)
	VALUES 
	($1, $2, $3, $4, $5) RETURNING id`

	row := tx.QueryRow(ctx, query, car.Make, car.Model, car.Trim, car.Year, car.VIN)
	var carId string
	if err := row.Scan(&carId); err != nil {
		return "", fmt.Errorf("failed to insert car: %w", err)
	}

	return carId, nil
}

func createUserCarRecord(ctx context.Context, tx pgx.Tx, userId, carId string) (string, error) {
	query := `
	INSERT INTO users_cars (user_id, car_id)
	VALUES 
	($1, $2) RETURNING id`

	row := tx.QueryRow(ctx, query, userId, carId)
	var userCarId string
	if err := row.Scan(&userCarId); err != nil {
		return "", fmt.Errorf("failed to insert user car: %w", err)
	}

	return userCarId, nil
}

func (s *Service) GetUsersCars(ctx context.Context, userId string) ([]Car, error) {
	if s.db == nil {
		return nil, ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(userId) == "" {
		return nil, ErrInvalidArg
	}

	query := `
	SELECT 
		c.id,
    	c.make,
    	c.model,
    	c.trim,
    	c.year,
    	c.vin,
    	c.created_at,
    	c.updated_at
	FROM cars c
	JOIN users_cars uc ON uc.car_id = c.id
	WHERE uc.user_id = $1`

	rows, err := s.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query for user cars: %w", err)
	}

	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var c Car
		if err := rows.Scan(c.id, c.Make, c.Model, c.Trim, c.Year, c.VIN, c.createdAt, c.updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		cars = append(cars, c)
	}

	return cars, nil
}
