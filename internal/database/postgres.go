package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/zercle/gofiber-skeleton/internal/config"
)

// PostgresDB wraps the database connection
type PostgresDB struct {
	*sql.DB
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	return db.DB.Close()
}

// HealthCheck performs a health check on the database
func (db *PostgresDB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetStats returns database connection pool statistics
func (db *PostgresDB) GetStats() sql.DBStats {
	return db.Stats()
}
