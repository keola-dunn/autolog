package car

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("the requested resource was not found")
)

type Car struct {
	id    string
	Make  string
	Model string
	Trim  string
	Year  int64
	VIN   string

	Color string

	// PublicId is the short 6 character ID assigned to each car for
	// ease of lookup. Eventual use case would be on a sticker, or as
	// a QR Code param.
	publicId string

	createdAt time.Time
	updatedAt time.Time
}

func (c *Car) Id() string {
	return c.id
}

func (c *Car) PublicId() string {
	return c.publicId
}

func (s *Service) CreateCar(ctx context.Context, userId string, car Car, nhtsaData NHTSAVPICData) error {
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

	publicId, _, err := s.generatePublicId(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate public id: %w", err)
	}

	car.publicId = publicId

	carId, err := createCarRecord(ctx, tx, car)
	if err != nil {
		return fmt.Errorf("failed to create car record: %w", err)
	}

	if _, err := createUserCarRecord(ctx, tx, userId, carId); err != nil {
		return fmt.Errorf("failed to create user car record: %w", err)
	}

	nhtsaData.carId = carId
	if err := createNHTSAVPICDataRecord(ctx, tx, nhtsaData); err != nil {
		return fmt.Errorf("failed to create nhtsa vpic data record: %w", err)
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
	INSERT INTO cars (make, model, trim, year, vin, color)
	VALUES 
	($1, $2, $3, $4, $5, $6) RETURNING id`

	row := tx.QueryRow(ctx, query, car.Make, car.Model, car.Trim, car.Year, car.VIN, car.Color)
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
		c.public_id,
    	c.make,
    	c.model,
    	c.trim,
    	c.year,
    	c.vin,
		c.color,
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
		if err := rows.Scan(&c.id, &c.publicId, &c.Make, &c.Model, &c.Trim,
			&c.Year, &c.VIN, &c.Color, &c.createdAt, &c.updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		cars = append(cars, c)
	}

	return cars, nil
}

type GetCarInput struct {
	VIN string

	PublicId string

	Id string

	// PlateNumber - DO NOT USE - edge cases to solve
	PlateNumber string
	// PlateState - DO NOT USE - edge cases to solve
	PlateState string
}

func (g *GetCarInput) valid() bool {
	return strings.TrimSpace(g.VIN) != "" ||
		strings.TrimSpace(g.PublicId) != "" ||
		(strings.TrimSpace(g.PlateNumber) != "" && strings.TrimSpace(g.PlateState) != "")
}

type GetCarOutput struct {
	Car
	Id        string
	PublicId  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Service) GetCar(ctx context.Context, input GetCarInput) (GetCarOutput, error) {
	if s.db == nil {
		return GetCarOutput{}, ErrMissingRequiredConfiguration
	}

	if !input.valid() {
		return GetCarOutput{}, ErrInvalidArg
	}

	var queryBuilder strings.Builder
	var conditionalQueryArgs = make([]string, 0, 3)
	var queryArgs []any

	// TODO: there's some license plate resolutions I need to make
	// ex. 1 plate could be applied to several cars
	// queryBuilder.WriteString(`
	// WITH latest_plates AS (
	// 	SELECT
	// 		l.*
	// 	FROM license_plates l
	// 	JOIN (
	// 		SELECT
	// 			l.car_id ,
	// 			MAX(l.created_at) max_created_at
	// 		FROM license_plates l
	// 		GROUP BY l.car_id
	// 	) pl ON pl.max_created_at = l.created_at
	// )
	// SELECT
	// 	c.id,
	// 	c.public_id,
	// 	c.make,
	// 	c.model,
	// 	c.trim,
	// 	c.YEAR,
	// 	c.vin,
	// 	c.created_at,
	// 	c.updated_at
	// FROM cars c
	// LEFT JOIN latest_plates lp ON lp.car_id = c.id WHERE `)

	queryBuilder.WriteString(`
	SELECT
		c.id,
		c.public_id,
		c.make,
		c.model,
		c.trim,
		c.year,
		c.vin,
		c.created_at,
		c.updated_at
	FROM cars c
	WHERE `)

	if strings.TrimSpace(input.VIN) != "" {
		queryArgs = append(queryArgs, strings.TrimSpace(input.VIN))
		conditionalQueryArgs = append(conditionalQueryArgs, fmt.Sprintf("c.vin = $%d", len(queryArgs)))
	}

	if strings.TrimSpace(input.PublicId) != "" {
		queryArgs = append(queryArgs, strings.TrimSpace(input.PublicId))
		conditionalQueryArgs = append(conditionalQueryArgs, fmt.Sprintf("c.public_id = $%d", len(queryArgs)))
	}

	if strings.TrimSpace(input.Id) != "" {
		queryArgs = append(queryArgs, strings.TrimSpace(input.Id))
		conditionalQueryArgs = append(conditionalQueryArgs, fmt.Sprintf("c.id = $%d", len(queryArgs)))
	}

	// if strings.TrimSpace(input.PlateNumber) != "" {
	// 	queryArgs = append(queryArgs, strings.TrimSpace(input.PlateNumber))
	// 	queryArgs = append(queryArgs, strings.TrimSpace(input.PlateState))
	// 	conditionalQueryArgs = append(conditionalQueryArgs, fmt.Sprintf("(lp.plate_number = $%d  AND lp.state = $%d)", len(conditionalQueryArgs)-1, len(conditionalQueryArgs)))
	// }

	queryBuilder.WriteString(strings.Join(conditionalQueryArgs, " OR "))

	var c Car
	row := s.db.QueryRow(ctx, queryBuilder.String(), queryArgs...)
	if err := row.Scan(&c.id, &c.publicId, &c.Make, &c.Model, &c.Trim,
		&c.Year, &c.VIN, &c.createdAt, &c.updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetCarOutput{}, ErrNotFound
		}
		return GetCarOutput{}, fmt.Errorf("failed to query for car: %w", err)
	}

	return GetCarOutput{
		Car:       c,
		Id:        c.id,
		PublicId:  c.publicId,
		CreatedAt: c.createdAt,
		UpdatedAt: c.updatedAt,
	}, nil
}
