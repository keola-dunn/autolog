package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrMissingRequiredConfiguration = errors.New("auth service is missing required configurations to perform this operation")

	ErrInvalidArg = errors.New("one or more of the provided arguments are invalid")
)

type CreateNewUserInput struct {
	Username string
	Email    string
	Password string
}

func (c *CreateNewUserInput) Valid() bool {
	// TODO: how could I scan/check for username creations that are offensive?
	if strings.TrimSpace(c.Email) == "" || strings.TrimSpace(c.Password) == "" || strings.TrimSpace(c.Email) == "" {
		return false
	}

	return true
}

// CreateNewUser creates a new user in the Auth service
// Takes a context and CreateNewUserInput as the inputs, returns UserID and an error as the output
func (s *Service) CreateNewUser(ctx context.Context, input CreateNewUserInput) (int64, error) {
	if s.db == nil {
		return 0, ErrMissingRequiredConfiguration
	}

	if !input.Valid() {
		return 0, ErrInvalidArg
	}

	salt := s.randomGenerator.RandomString(s.saltLength)

	passwordHash := argon2.IDKey([]byte(input.Password), []byte(salt), 1, 64*1024, 4, 32)

	query := `
	INSERT INTO users(username, salt, password_hash, email)
	VALUES ($1, $2, $3, $4) RETURNING id`

	row := s.db.QueryRow(ctx, query, input.Username, salt, passwordHash, input.Email)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to exec create new user query: %w", err)
	}

	return id, nil
}
