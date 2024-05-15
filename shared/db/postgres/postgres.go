package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)
	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User, cfg.Password, hostPort, cfg.Database, cfg.SSLMode))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
