package shop

import (
	"context"
	"fmt"
	"strings"
)

type Role string

const (
	RoleEmployee = Role("employee")

	RoleOwner = Role("owner")
)

type CreateEmployeeInput struct {
	ShopId          string
	UserId          string
	Role            Role
	CreatedByUserId string
}

func (s *Service) CreateEmployee(ctx context.Context, input CreateEmployeeInput) (string, error) {
	if s.db == nil {
		return "", ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(input.ShopId) == "" ||
		strings.TrimSpace(input.UserId) == "" ||
		strings.TrimSpace(string(input.Role)) == "" ||
		strings.TrimSpace(input.CreatedByUserId) == "" {
		return "", ErrInvalidArg
	}

	query := `
	INSERT INTO employees (user_id, shop_id, "role", created_by) 
	VALUES ($1, $2, $3, $4) RETURNING id`

	row := s.db.QueryRow(ctx, query, strings.TrimSpace(input.UserId),
		strings.TrimSpace(input.ShopId), input.Role, strings.TrimSpace(input.CreatedByUserId))

	var employeeId string
	if err := row.Scan(&employeeId); err != nil {
		return "", fmt.Errorf("failed to insert employee: %w", err)
	}

	return employeeId, nil
}
