package car

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type ServiceLog struct {
	Type    string
	Date    time.Time
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
	INSERT INTO service_logs (user_id, car_id, type, date, details, notes)
	VALUES
	($1, $2, $3, $4, $5, $6) RETURNING id`

	row := s.db.QueryRow(ctx, query, userId, carId, serviceLog.Type, serviceLog.Date, serviceLog.Details, serviceLog.Notes)

	var serviceLogId string
	if err := row.Scan(&serviceLogId); err != nil {
		return "", fmt.Errorf("failed to scan service log row: %w", err)
	}

	return serviceLogId, nil
}
