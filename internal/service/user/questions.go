package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type SecurityQuestion struct {
	Id        string
	Question  string
	CreatedAt time.Time
}

func (s *Service) GetSecurityQuestions(ctx context.Context) ([]SecurityQuestion, error) {
	query := `
	SELECT 
		id,
    	question,
    	created_at
	FROM security_questions`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query for questions: %w", err)
	}
	defer rows.Close()

	var questions = []SecurityQuestion{}

	for rows.Next() {
		var q SecurityQuestion

		if err := rows.Scan(
			&q.Id,
			&q.Question,
			&q.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan question row as expected: %w", err)
		}

		questions = append(questions, q)
	}

	return questions, nil
}

type UserSecurityQuestion struct {
	QuestionId string
	Answer     string
}

type userSecurityQuestionRecord struct {
	id         string
	userId     string
	questionId string
	answerHash string
	salt       string
	createdAt  time.Time
}

func createUserSecurityQuestions(ctx context.Context, dbTransaction pgx.Tx, questions []userSecurityQuestionRecord) error {
	var query strings.Builder
	query.WriteString(`
	INSERT INTO users_security_questions (user_id, question_id, answer_hash, salt)
	VALUES `)

	var args = []any{}

	for i, q := range questions {
		args = append(args, q.userId, q.questionId, q.answerHash, q.salt)

		query.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)",
			len(args)-3, len(args)-2, len(args)-1, len(args),
		))

		if i != len(questions)-1 {
			query.WriteString(",\n")
		}
	}

	if _, err := dbTransaction.Exec(ctx, query.String(), args...); err != nil {
		return fmt.Errorf("failed to insert user security questions: %w", err)
	}

	return nil
}
