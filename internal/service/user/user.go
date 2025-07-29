package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	CreateNewUser(context.Context, CreateNewUserInput) (string, error)
	ValidateCredentials(ctx context.Context, user, password string) (bool, string, error)

	DoesUsernameOrEmailExist(ctx context.Context, username, email string) (bool, bool, error)

	GetSecurityQuestions(context.Context) ([]SecurityQuestion, error)

	GetUserRole(ctx context.Context, userId string) (GetUserRoleOutput, error)
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
	Username          string
	Email             string
	Password          string
	SecurityQuestions []UserSecurityQuestion
	Role              Role
}

func (c *CreateNewUserInput) Valid() bool {
	// TODO: how could I scan/check for username creations that are offensive?
	if strings.TrimSpace(c.Email) == "" ||
		strings.TrimSpace(c.Password) == "" ||
		strings.TrimSpace(c.Email) == "" ||
		len(c.SecurityQuestions) < 3 {
		return false
	}

	return true
}

// CreateNewUser creates a new user in the Auth service
// Takes a context and CreateNewUserInput as the inputs, returns UserID and an error as the output
func (s *Service) CreateNewUser(ctx context.Context, input CreateNewUserInput) (string, error) {
	if s.db == nil {
		return "", ErrMissingRequiredConfiguration
	}

	if !input.Valid() {
		return "", ErrInvalidArg
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	salt := s.randomGenerator.RandomString(s.saltLength)

	passwordHash := s.passwordHash(input.Password, salt)

	userId, err := createNewUserRecord(ctx, tx, input.Username, salt, string(passwordHash), input.Email)
	if err != nil {
		// TODO: figure out the error when a unique constraint on username or email is violated
		return "", fmt.Errorf("failed to create new user record: %w", err)
	}

	userSecQuestions := make([]userSecurityQuestionRecord, 0, len(input.SecurityQuestions))
	for _, question := range input.SecurityQuestions {
		salt = s.randomGenerator.RandomString(s.saltLength)
		answerHash := s.passwordHash(question.Answer, salt)

		userSecQuestions = append(userSecQuestions, userSecurityQuestionRecord{
			questionId: question.QuestionId,
			answerHash: string(answerHash),
			salt:       salt,
			userId:     userId,
		})
	}

	if err := createUserSecurityQuestions(ctx, tx, userSecQuestions); err != nil {
		return "", fmt.Errorf("failed to create user security questions: %w", err)
	}

	if err := createUserRoleRecord(ctx, tx, userId, input.Role); err != nil {
		return "", fmt.Errorf("failed to create user role record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit transaction to db: %w", err)
	}

	return userId, nil
}

func createNewUserRecord(ctx context.Context, dbTransaction pgx.Tx, username, salt, passwordHash, email string) (string, error) {
	query := `
	INSERT INTO users(username, salt, password_hash, email)
	VALUES ($1, $2, $3, $4) RETURNING id`

	row := dbTransaction.QueryRow(ctx, query, username, salt, passwordHash, email)

	var id pgtype.UUID
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("failed to insert new user: %w", err)
	}

	return id.String(), nil
}

func (s *Service) passwordHash(password, salt string) string {
	pwHash := argon2.Key([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(pwHash)
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

	var userId pgtype.UUID
	var salt string
	var storedPasswordHash string
	if err := row.Scan(&userId, &salt, &storedPasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", nil
		}

		return false, "", fmt.Errorf("failed to query for valid user credentials: %w", err)
	}

	providedHash := s.passwordHash(password, salt)
	if providedHash == storedPasswordHash {
		return true, userId.String(), nil
	}

	return false, "", nil
}

// DoesUsernameOrEmailExist checks if a provided username or email exists in the users table already.
// Both a username and an email must be provided. This is intended to check both.
// returns bool (username exists), bool (email exists), and an error
func (s *Service) DoesUsernameOrEmailExist(ctx context.Context, username, email string) (bool, bool, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(email) == "" {
		return false, false, ErrInvalidArg
	}

	query := `
	WITH existing_user_records AS (
		SELECT 
			username, 
			email
		FROM users 
		WHERE username = $1 OR email = $2
	)
	SELECT 
		CASE 
			WHEN (SELECT true FROM existing_user_records WHERE username = $1) THEN true
		ELSE false
		END as username_exists,
		CASE 
			WHEN (SELECT true FROM existing_user_records WHERE email = $2) THEN true
		ELSE false
		END as email_exists
	`
	row := s.db.QueryRow(ctx, query, strings.TrimSpace(username), strings.TrimSpace(email))

	var usernameExists bool
	var emailExists bool
	if err := row.Scan(&usernameExists, &emailExists); err != nil {
		return usernameExists, emailExists, fmt.Errorf("failed to query if username or email exists: %w", err)
	}

	return usernameExists, emailExists, nil
}
