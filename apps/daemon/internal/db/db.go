package db

import (
	"context"
	"fmt"

	"github.com/enki/daemon/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPoolFromConfig creates a new PostgreSQL connection pool from configuration
func NewPoolFromConfig(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dbURL := cfg.DatabaseURL()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
