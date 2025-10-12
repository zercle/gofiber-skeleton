package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samber/do/v2"
)

type Database struct {
	DB     *sqlx.DB
	Pool   *pgxpool.Pool
	config Config
}

type Config interface {
	GetDatabaseDSN() string
}

// NewDatabase creates a new database connection
func NewDatabase(injector do.Injector) (*Database, error) {
	config := do.MustInvoke[Config](injector)

	// Create pgx pool for SQLC
	pool, err := pgxpool.New(context.Background(), config.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Create sqlx DB for compatibility
	db, err := sqlx.Connect("postgres", config.GetDatabaseDSN())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	return &Database{
		DB:     db,
		Pool:   pool,
		config: config,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	var errs []error

	if d.Pool != nil {
		d.Pool.Close()
	}

	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close sqlx db: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple errors occurred during close: %v", errs)
	}

	return nil
}

// Health checks the database connection health
func (d *Database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetDB returns the underlying sqlx DB instance
func (d *Database) GetDB() *sqlx.DB {
	return d.DB
}

// GetPool returns the underlying pgx pool for SQLC
func (d *Database) GetPool() *pgxpool.Pool {
	return d.Pool
}

// BeginTx starts a new transaction
func (d *Database) BeginTx() (*sqlx.Tx, error) {
	return d.DB.Beginx()
}

// NewTestDB creates a database connection for testing
func NewTestDB(dsn string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Configure connection pool for testing
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(1 * time.Minute)

	return &Database{
		DB: db,
	}, nil
}
