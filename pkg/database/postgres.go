package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDB creates a new connection pool to the database
func NewPostgresDB(dsn string) (*pgxpool.Pool, error) {
	// 1. Parse the configuration string
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// 2. Set connection settings (optional but recommended for production)
	config.MaxConns = 10                   // Max simultaneous connections
	config.MinConns = 2                    // Min idle connections
	config.MaxConnLifetime = 1 * time.Hour // Recycle connections every hour

	// 3. Create the connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to creare connection pool: %w", err)
	}

	return pool, nil
}
