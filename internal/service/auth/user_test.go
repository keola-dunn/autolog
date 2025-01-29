package auth_test

import (
	"context"
	"errors"

	"testing"

	"github.com/keola-dunn/autolog/internal/random"
	"github.com/keola-dunn/autolog/internal/service/auth"
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
		input auth.CreateNewUserInput

		expectedQuery     string
		expectedQueryArgs []interface{}
		fakeQueryRows     *pgxmock.Rows
		fakeQueryErr      error

		expectedUserId int64
		expectedErr    error
	}{
		{
			name:              "InvalidArg",
			input:             auth.CreateNewUserInput{},
			expectedQuery:     "",
			expectedQueryArgs: nil,
			fakeQueryRows:     nil,
			fakeQueryErr:      nil,
			expectedUserId:    0,
			expectedErr:       auth.ErrInvalidArg,
		},
		{
			name: "DbError",
			input: auth.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
			},
			expectedQuery: "",
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
			name: "DbError",
			input: auth.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
			},
			expectedQuery: "",
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
			input: auth.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
			},
			expectedQuery: "",
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
			db, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("failed to create new test postgres db: %v", err)
			}

			db.ExpectQuery(test.expectedQuery).
				WithArgs(test.expectedQueryArgs...).
				WillReturnRows(test.fakeQueryRows).
				WillReturnError(test.fakeQueryErr)

			service := auth.NewService(auth.ServiceConfig{
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
