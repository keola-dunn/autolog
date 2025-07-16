package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/keola-dunn/autolog/internal/platform/postgres"
	"github.com/keola-dunn/autolog/internal/random"
	"golang.org/x/crypto/argon2"
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

	// SaltLength sets the length of the password salts generated for the auth service.
	// Defaults to 8.
	SaltLength int64
}

type ServiceIface interface {
	CreateNewUser(context.Context, CreateNewUserInput) (int64, error)
	ValidateCredentials(ctx context.Context, user, password string) (bool, string, error)
}

type Service struct {
	db              postgres.ConnectionPool
	randomGenerator random.ServiceIface
	saltLength      int64
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.RandomGenerator == nil {
		cfg.RandomGenerator = random.NewService()
	}

	if cfg.SaltLength <= 0 {
		cfg.SaltLength = 8
	}

	return &Service{
		db:              cfg.DB,
		randomGenerator: cfg.RandomGenerator,
		saltLength:      cfg.SaltLength,
	}
}

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

	passwordHash := s.passwordHash(input.Password, salt)

	query := `
	INSERT INTO users(username, salt, password_hash, email)
	VALUES ($1, $2, $3, $4) RETURNING id`

	row := s.db.QueryRow(ctx, query, input.Username, salt, string(passwordHash), input.Email)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to exec create new user query: %w", err)
	}
	// TODO: figure out the error when a unique constraint on username or email is violated
	return id, nil
}

func (s *Service) passwordHash(password, salt string) []byte {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
}

// ValidateCredentials will check the provided credentials against the database. This
// is meant to be used as a login method. Returns
// true if the credentials are good, false otherwise, the user id
// (if valid) and an error.
func (s *Service) ValidateCredentials(ctx context.Context, user, password string) (bool, string, error) {
	if s.db == nil {
		return false, "", ErrMissingRequiredConfiguration
	}

	if strings.TrimSpace(user) == "" || strings.TrimSpace(password) == "" {
		return false, "", ErrInvalidArg
	}

	query := `
	SELECT 
		u.id,
		u.salt, 
		u.password_hash
	FROM users u
	WHERE u.username = $1 OR u.email = $1`

	row := s.db.QueryRow(ctx, query, user)

	var userId string
	var salt string
	var storedPasswordHash string
	if err := row.Scan(&userId, &salt, &storedPasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", nil
		}

		return false, "", fmt.Errorf("failed to query for valid user credentials: %w", err)
	}

	providedHash := s.passwordHash(password, salt)
	if bytes.Equal(providedHash, []byte(storedPasswordHash)) {
		return true, userId, nil
	}

	return false, "", nil
}
