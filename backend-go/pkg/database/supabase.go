package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// NewSupabaseDB creates a new database connection to Supabase
func NewSupabaseDB(connectionString string) (*DB, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("database connection string is empty")
	}

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Connection pool settings optimized for API server
	config.MaxConns = 25              // Maximum number of connections in the pool
	config.MinConns = 5               // Minimum number of connections to maintain
	config.MaxConnLifetime = 0        // Connections live forever
	config.MaxConnIdleTime = 0        // Idle connections are not closed
	config.HealthCheckPeriod = 0      // Disable health check (Supabase handles this)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Health checks if the database connection is alive
func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
