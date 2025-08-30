package car

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Service) generatePublicId(ctx context.Context) (string, int64, error) {
	var attempts int64
	for {
		attempts++
		id := s.randomGenerator.RandomUpperAlphanumericString(s.publicIdLength)

		checkIdQuery := `
		SELECT 
			1
		FROM cars c 
		WHERE c.public_id = $1`

		row := s.db.QueryRow(ctx, checkIdQuery, id)
		var exists bool

		if err := row.Scan(&exists); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", attempts, fmt.Errorf("failed to check if id %s exists: %w", id, err)
		}

		if !exists {
			return id, attempts, nil
		}
	}
}
