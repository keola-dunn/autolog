package user_test

import (
	"context"
	"errors"

	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/user"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

type fakeRandomService struct {
	random.ServiceIface
}

func (f *fakeRandomService) RandomString(_ int64) string {
	return "fakerandomstring"
}

func TestCreateNewUser(t *testing.T) {
	tests := []struct {
		name  string
		input user.CreateNewUserInput

		expectedQuery     string
		expectedQueryArgs []interface{}
		fakeQueryRows     *pgxmock.Rows
		fakeQueryErr      error

		expectedUserId int64
		expectedErr    error
	}{
		{
			name:              "InvalidArg",
			input:             user.CreateNewUserInput{},
			expectedQuery:     "",
			expectedQueryArgs: nil,
			fakeQueryRows:     nil,
			fakeQueryErr:      nil,
			expectedUserId:    0,
			expectedErr:       user.ErrInvalidArg,
		},
		{
			name: "DbError",
			input: user.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
			},
			expectedQuery: "INSERT INTO users(username, salt, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id",
			expectedQueryArgs: []interface{}{
				"TestUsername",
				"fakerandomstring",
				string([]uint8{25, 57, 194, 158, 249, 111, 195, 48, 173, 95, 240, 160, 80, 177, 217, 217, 70, 114, 217, 170, 140, 12, 47, 250, 45, 133, 246, 213, 54, 209, 213, 175}),
				"TestEmail"},
			fakeQueryRows:  pgxmock.NewRows(nil),
			fakeQueryErr:   errors.New("fake db error"),
			expectedUserId: 0,
			expectedErr:    errors.New("failed to exec create new user query: fake db error"),
		},
		{
			name: "Success",
			input: user.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
			},
			expectedQuery: "INSERT INTO users(username, salt, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id",
			expectedQueryArgs: []interface{}{
				"TestUsername",
				"fakerandomstring",
				string([]uint8{25, 57, 194, 158, 249, 111, 195, 48, 173, 95, 240, 160, 80, 177, 217, 217, 70, 114, 217, 170, 140, 12, 47, 250, 45, 133, 246, 213, 54, 209, 213, 175}),
				"TestEmail"},
			fakeQueryRows: pgxmock.NewRows([]string{
				"id",
			}).AddRow(int64(1001)),
			fakeQueryErr:   nil,
			expectedUserId: 1001,
			expectedErr:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("failed to create new test postgres db: %v", err)
			}
			defer db.Close(context.Background())

			db.ExpectQuery(test.expectedQuery).
				WithArgs(test.expectedQueryArgs...).
				WillReturnRows(test.fakeQueryRows).
				WillReturnError(test.fakeQueryErr)

			service := user.NewService(user.ServiceConfig{
				DB:              db,
				RandomGenerator: &fakeRandomService{},
			})

			userId, err := service.CreateNewUser(context.TODO(), test.input)
			if err != test.expectedErr && (err == nil || test.expectedErr == nil || err.Error() != test.expectedErr.Error()) {
				t.Errorf("expected error:\n%v\ndoes not match actual:\n%v", test.expectedErr, err)
			}
			require.Equal(t, test.expectedUserId, userId, "userIdComparison")
		})
	}
}

func TestValidateCredentials(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string

		expectedQuery     string
		expectedQueryArgs []interface{}
		fakeQueryRows     *pgxmock.Rows
		fakeQueryErr      error

		expectedValid  bool
		expectedUserId string
		expectedErr    error
	}{
		{
			name:              "InvalidArg",
			user:              "",
			password:          "",
			expectedQuery:     "",
			expectedQueryArgs: nil,
			fakeQueryRows:     nil,
			fakeQueryErr:      nil,
			expectedValid:     false,
			expectedErr:       user.ErrInvalidArg,
		},
		{
			name:              "DbError",
			user:              "Username",
			password:          "Password",
			expectedQuery:     "SELECT u.id, u.salt, u.password_hash FROM users u WHERE u.username = $1 OR u.email = $1",
			expectedQueryArgs: []interface{}{"Username"},
			fakeQueryRows:     pgxmock.NewRows(nil),
			fakeQueryErr:      errors.New("fake db error"),
			expectedValid:     false,
			expectedErr:       errors.New("failed to query for valid user credentials: fake db error"),
		},
		{
			name:              "NoRows",
			user:              "Username",
			password:          "Password",
			expectedQuery:     "SELECT u.id, u.salt, u.password_hash FROM users u WHERE u.username = $1 OR u.email = $1",
			expectedQueryArgs: []interface{}{"Username"},
			fakeQueryRows:     pgxmock.NewRows(nil),
			fakeQueryErr:      pgx.ErrNoRows,
			expectedValid:     false,
			expectedErr:       nil,
		},
		{
			name:              "Invalid",
			user:              "Username",
			password:          "Password",
			expectedQuery:     "SELECT u.id, u.salt, u.password_hash FROM users u WHERE u.username = $1 OR u.email = $1",
			expectedQueryArgs: []interface{}{"Username"},
			fakeQueryRows: pgxmock.NewRows([]string{"id", "salt", "password_hash"}).AddRows(
				[]interface{}{
					"UserID",
					"fakerandomstring",
					string([]uint8{25, 57, 194, 158, 249, 111, 195, 48, 173, 95, 240, 160, 80, 177, 217, 217, 70, 114, 217, 170, 140, 12, 47, 250, 45, 133, 246, 213, 54, 209, 213, 175})},
			),
			fakeQueryErr:   nil,
			expectedValid:  false,
			expectedUserId: "",
			expectedErr:    nil,
		},
		{
			name:              "ValidCredentials",
			user:              "Username",
			password:          "TestPassword",
			expectedQuery:     "SELECT u.id, u.salt, u.password_hash FROM users u WHERE u.username = $1 OR u.email = $1",
			expectedQueryArgs: []interface{}{"Username"},
			fakeQueryRows: pgxmock.NewRows([]string{"id", "salt", "password_hash"}).AddRows(
				[]interface{}{
					"UserID",
					"fakerandomstring",
					string([]uint8{25, 57, 194, 158, 249, 111, 195, 48, 173, 95, 240, 160, 80, 177, 217, 217, 70, 114, 217, 170, 140, 12, 47, 250, 45, 133, 246, 213, 54, 209, 213, 175})},
			),
			fakeQueryErr:   nil,
			expectedValid:  true,
			expectedUserId: "UserID",
			expectedErr:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("failed to create new test postgres db: %v", err)
			}
			defer db.Close(context.Background())

			db.ExpectQuery(test.expectedQuery).
				WithArgs(test.expectedQueryArgs...).
				WillReturnRows(test.fakeQueryRows).
				WillReturnError(test.fakeQueryErr)

			service := user.NewService(user.ServiceConfig{
				DB:              db,
				RandomGenerator: &fakeRandomService{},
			})

			valid, userId, err := service.ValidateCredentials(context.TODO(), test.user, test.password)
			if err != test.expectedErr && (err == nil || test.expectedErr == nil || err.Error() != test.expectedErr.Error()) {
				t.Errorf("expected error:\n%v\ndoes not match actual:\n%v", test.expectedErr, err)
			}
			require.Equal(t, test.expectedValid, valid, "validComparison")
			require.Equal(t, test.expectedUserId, userId, "userIdComparison")
		})
	}
}
