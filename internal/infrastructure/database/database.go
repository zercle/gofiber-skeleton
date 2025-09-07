package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

// Database represents the database connection pool.
type Database struct {
	Pool *pgxpool.Pool
}

// NewConnection creates a new database connection pool.
func NewDB(cfg *config.Config) (*Database, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse database URL and configure connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	poolConfig.MaxConnLifetime = time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = 30 * time.Second // Add idle time to match common practices

	// Establish connection
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cleanup := func() {
		pool.Close()
	}

	return &Database{Pool: pool}, cleanup, nil
}

// Close closes the database connection pool.
func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

// Health checks the database connection.
func (d *Database) Health(ctx context.Context) error {
	return d.Pool.Ping(ctx)
}
