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
	testUserId := "e186aa27-10d4-4f06-907f-ec1a37174a98"
	tests := []struct {
		name  string
		input user.CreateNewUserInput

		dbFunc         func(db pgxmock.PgxConnIface)
		expectedUserId string
		expectedErr    error
	}{
		{
			name:           "InvalidArg",
			input:          user.CreateNewUserInput{},
			dbFunc:         func(db pgxmock.PgxConnIface) {},
			expectedUserId: "",
			expectedErr:    user.ErrInvalidArg,
		},
		{
			name: "DbError-CreateUserRecord",
			input: user.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
				SecurityQuestions: []user.UserSecurityQuestion{
					{}, {}, {},
				},
			},
			dbFunc: func(db pgxmock.PgxConnIface) {
				db.ExpectBegin()

				db.ExpectQuery("INSERT INTO users(username, salt, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id").
					WithArgs(
						"TestUsername",
						"fakerandomstring",
						"m2LHi/PGOgAmCn17BQx8wTp9JZdc8lCBELH2NPsvSVs",
						"TestEmail").
					WillReturnRows(pgxmock.NewRows(nil)).
					WillReturnError(errors.New("fake db error"))

				db.ExpectRollback()
			},
			expectedUserId: "",
			expectedErr:    errors.New("failed to create new user record: failed to insert new user: fake db error"),
		},
		{
			name: "DbError-CreateUserSecurityQuestionRecords",
			input: user.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
				SecurityQuestions: []user.UserSecurityQuestion{
					{
						QuestionId: "TestQuestionId",
						Answer:     "Test Answer 1",
					}, {
						QuestionId: "TestQuestionId2",
						Answer:     "Test Answer 2",
					}, {
						QuestionId: "TestQuestionId3",
						Answer:     "Test Answer 3",
					},
				},
			},
			dbFunc: func(db pgxmock.PgxConnIface) {
				db.ExpectBegin()

				db.ExpectQuery("INSERT INTO users(username, salt, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id").
					WithArgs(
						"TestUsername",
						"fakerandomstring",
						"m2LHi/PGOgAmCn17BQx8wTp9JZdc8lCBELH2NPsvSVs",
						"TestEmail").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(testUserId))

				db.ExpectExec((`
				INSERT INTO users_security_questions (user_id, question_id, answer_hash, salt) 
				VALUES 
					($1, $2, $3, $4), 
					($5, $6, $7, $8), 
					($9, $10, $11, $12)`)).
					WithArgs(
						testUserId, "TestQuestionId", "QxPGM9/ix/57HRCs3ZbIqaVI7nbKmwcZHJIeLIgtc5Y", "fakerandomstring",
						testUserId, "TestQuestionId2", "b+hYERj3RxnD8+ijNE3Ot6yS+q62VNGlXvJnWZ84o5U", "fakerandomstring",
						testUserId, "TestQuestionId3", "zvdpXECZZTj4uu3A8Tleo9RKRMUx+Wz9Xqbt+8Evug8", "fakerandomstring").
					WillReturnError(errors.New("fake db error"))

				db.ExpectRollback()
			},
			expectedUserId: "",
			expectedErr:    errors.New("failed to create user security questions: failed to insert user security questions: fake db error"),
		},
		{
			name: "Success",
			input: user.CreateNewUserInput{
				Username: "TestUsername",
				Email:    "TestEmail",
				Password: "TestPassword",
				SecurityQuestions: []user.UserSecurityQuestion{
					{
						QuestionId: "TestQuestionId",
						Answer:     "Test Answer 1",
					}, {
						QuestionId: "TestQuestionId2",
						Answer:     "Test Answer 2",
					}, {
						QuestionId: "TestQuestionId3",
						Answer:     "Test Answer 3",
					},
				},
			},
			dbFunc: func(db pgxmock.PgxConnIface) {
				db.ExpectBegin()

				db.ExpectQuery("INSERT INTO users(username, salt, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id").
					WithArgs(
						"TestUsername",
						"fakerandomstring",
						"m2LHi/PGOgAmCn17BQx8wTp9JZdc8lCBELH2NPsvSVs",
						"TestEmail").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(testUserId))

				db.ExpectExec((`
				INSERT INTO users_security_questions (user_id, question_id, answer_hash, salt) 
				VALUES 
					($1, $2, $3, $4), 
					($5, $6, $7, $8), 
					($9, $10, $11, $12)`)).
					WithArgs(
						testUserId, "TestQuestionId", "QxPGM9/ix/57HRCs3ZbIqaVI7nbKmwcZHJIeLIgtc5Y", "fakerandomstring",
						testUserId, "TestQuestionId2", "b+hYERj3RxnD8+ijNE3Ot6yS+q62VNGlXvJnWZ84o5U", "fakerandomstring",
						testUserId, "TestQuestionId3", "zvdpXECZZTj4uu3A8Tleo9RKRMUx+Wz9Xqbt+8Evug8", "fakerandomstring").
					WillReturnResult(pgxmock.NewResult("insert", 3))

				db.ExpectCommit()
			},
			expectedUserId: testUserId,
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

			test.dbFunc(db)

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
	testUserId := "e186aa27-10d4-4f06-907f-ec1a37174a98"
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
					testUserId,
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
					testUserId,
					"fakerandomstring",
					"m2LHi/PGOgAmCn17BQx8wTp9JZdc8lCBELH2NPsvSVs"},
			),
			fakeQueryErr:   nil,
			expectedValid:  true,
			expectedUserId: testUserId,
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
