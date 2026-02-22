package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func RunMigrations(dsn string) {
	var m *migrate.Migrate
	var err error

	for i := range 5 {
		m, err = migrate.New("file://migrations", dsn)
		if err == nil {
			break
		}
		log.Printf("Waiting for database to be ready for migrations... (attempt %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Migration initialization failed: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ Database schema is up to date.")
		} else {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	} else {
		log.Println("🚀 Database migrated successfully!")
	}
}
