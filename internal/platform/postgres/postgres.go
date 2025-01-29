package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Ping(context.Context) error
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type ConnectionConfig struct {
	User string

	Password string

	DBHost string

	Port int64

	DBName string

	SSLMode string

	Schema string
}

func (c *ConnectionConfig) connectionString() string {
	if strings.TrimSpace(c.SSLMode) == "" {
		c.SSLMode = "verify-ca"
	}
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s search_path=%s",
		c.User, c.Password, c.DBHost, c.Port, c.DBName, c.SSLMode, c.Schema)
}

func NewConnectionPool(ctx context.Context, cfg ConnectionConfig) (ConnectionPool, error) {
	connpool, err := pgxpool.New(ctx, cfg.connectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to create new pool: %w", err)
	}

	return connpool, nil
}
