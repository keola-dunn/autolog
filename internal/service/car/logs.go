package car

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ServiceLog struct {
	Type    string
	Date    time.Time
	Mileage int64
	Details any
	Notes   string

	createdAt  time.Time
	updated_at time.Time
}

func (s *Service) CreateServiceLog(ctx context.Context, serviceLog ServiceLog, userId, carId string) (string, error) {
	if s.db == nil {
		return "", ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(userId) == "" ||
		strings.TrimSpace(carId) == "" {
		return "", ErrInvalidArg
	}

	query := `
	INSERT INTO service_logs (user_id, car_id, type, date, mileage, details, notes)
	VALUES
	($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	row := s.db.QueryRow(ctx, query, userId, carId, serviceLog.Type, serviceLog.Date, serviceLog.Mileage, serviceLog.Details, serviceLog.Notes)

	var serviceLogId string
	if err := row.Scan(&serviceLogId); err != nil {
		return "", fmt.Errorf("failed to scan service log row: %w", err)
	}

	return serviceLogId, nil
}

type ServiceLogSummary struct {
	Services map[string]struct {
		Count              int
		LastService        time.Time
		LastServiceMileage int64
	}
}

func (s *Service) GetServiceLogSummary(ctx context.Context, carId string) (ServiceLogSummary, error) {
	if strings.TrimSpace(carId) == "" {
		return ServiceLogSummary{}, ErrInvalidArg
	}

	query := `
	SELECT
		sl.type,
		sl.date,
		sl.mileage
	FROM service_logs sl
	WHERE 
		sl.car_id = $1`

	rows, err := s.db.Query(ctx, query, strings.TrimSpace(carId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ServiceLogSummary{}, ErrNotFound
		}
		return ServiceLogSummary{}, fmt.Errorf("failed to query for service logs: %w", err)
	}

	defer rows.Close()

	var output = ServiceLogSummary{
		Services: make(map[string]struct {
			Count              int
			LastService        time.Time
			LastServiceMileage int64
		}),
	}

	for rows.Next() {
		var serviceType string
		var serviceDate time.Time
		var mileage int64
		if err := rows.Scan(&serviceType, &serviceDate, &mileage); err != nil {
			return output, fmt.Errorf("failed to scan service log row as expected: %w", err)
		}

		serviceRecords, ok := output.Services[serviceType]
		if !ok {
			serviceRecords = struct {
				Count              int
				LastService        time.Time
				LastServiceMileage int64
			}{
				Count:              1,
				LastService:        serviceDate,
				LastServiceMileage: mileage,
			}
		} else {
			serviceRecords.Count++
			if serviceRecords.LastService.Before(serviceDate) {
				serviceRecords.LastService = serviceDate
				serviceRecords.LastServiceMileage = mileage
			}
		}

		output.Services[serviceType] = serviceRecords
	}

	return output, nil
}
