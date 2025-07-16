package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	Host string

	Port int64

	DBName string

	SSLMode string

	Schema string
}

type ConnectionPoolConfig struct {
	ConnectionConfig

	MaxConnections        int32
	MinConnections        int32
	MaxConnectionIdleTime time.Duration
}

func (c *ConnectionConfig) connectionString() string {
	// if strings.TrimSpace(c.SSLMode) == "" {
	// 	c.SSLMode = "verify-ca"
	// }

	var connStr strings.Builder
	connStr.WriteString(fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.SSLMode))

	if strings.TrimSpace(c.SSLMode) != "" {
		connStr.WriteString(fmt.Sprintf(" sslmode=%s", c.SSLMode))
	}

	if strings.TrimSpace(c.DBName) != "" {
		connStr.WriteString(fmt.Sprintf(" dbname=%s", c.DBName))
	}

	if strings.TrimSpace(c.Schema) != "" {
		connStr.WriteString(fmt.Sprintf(" search_path=%s", c.Schema))
	}

	return connStr.String()
}

func NewConnectionPool(ctx context.Context, cfg ConnectionPoolConfig) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.connectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse db connection config: %w", err)
	}

	dbConfig.MaxConns = cfg.MaxConnections
	dbConfig.MinConns = cfg.MinConnections
	dbConfig.MaxConnIdleTime = cfg.MaxConnectionIdleTime

	connpool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pool: %w", err)
	}

	return connpool, nil
}
