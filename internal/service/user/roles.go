package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Role string

const (
	// admin is the role that will allow full data access
	RoleAdmin = Role("admin")

	// user is the default role for users
	RoleUser = Role("user")
)

type GetUserRoleOutput struct {
	id        string
	userId    string
	Role      Role
	createdAt time.Time
	updatedAt time.Time
}

func (s *Service) GetUserRole(ctx context.Context, userId string) (GetUserRoleOutput, error) {
	if s.db == nil {
		return GetUserRoleOutput{}, ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(userId) == "" {
		return GetUserRoleOutput{}, ErrInvalidArg
	}

	query := `
	SELECT 
		ur.id,
    	ur.user_id,
    	ur.created_at,
 		ur.updated_at,
		r.role
	FROM users_roles ur
	JOIN roles r ON r.id = ur.role_id
	WHERE 
		ur.user_id = $1
	ORDER BY ur.created_at DESC 
	LIMIT 1`

	var output GetUserRoleOutput
	row := s.db.QueryRow(ctx, query, userId)
	if err := row.Scan(&output.id, &output.userId,
		&output.createdAt,
		&output.updatedAt,
		&output.Role); err != nil {
		return GetUserRoleOutput{}, fmt.Errorf("failed to query for user role: %w", err)
	}

	return output, nil
}

// createUserRoleRecord creates a new record in the users_roles table. This decision was made
// to keep a record of any changes to user roles.
func createUserRoleRecord(ctx context.Context, tx pgx.Tx, userId string, role Role) error {
	if strings.TrimSpace(userId) == "" ||
		strings.TrimSpace(string(role)) == "" ||
		tx == nil {
		return ErrInvalidArg
	}

	query := `
	WITH role AS (
		SELECT 
			id
		FROM roles r
		WHERE r.role = $1
		LIMIT 1
	)
	INSERT INTO users_roles (user_id, role_id) 
	VALUES 
	($2, (SELECT id FROM role))`

	if _, err := tx.Exec(ctx, query, string(role), userId); err != nil {
		return fmt.Errorf("failed to insert user role record: %w", err)
	}

	return nil
}
