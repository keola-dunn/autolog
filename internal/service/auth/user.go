package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingRequiredConfiguration = errors.New("auth service is missing required configurations to perform this operation")

	ErrInvalidArg = errors.New("one or more of the provided arguments are invalid")
)

type CreateNewUserInput struct {
	Username string
	Email    string
	Password string

	// how much info do I want to collect from the get go? Probably as little as possible
	ZipCode     string
	PhoneNumber string
	FirstName   string
	LastName    string
	Suffix      string
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

	query := ``

	row := s.db.QueryRow(ctx, query)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to exec create new user query: %w", err)
	}

	return id, nil
}
