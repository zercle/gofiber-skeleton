package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Pool *pgxpool.Pool
	DB   *sql.DB
}

func NewDatabase(ctx context.Context, dsn string, maxOpenConns, maxIdleConns int) (*Database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	return &Database{
		Pool: pool,
		DB:   db,
	}, nil
}

func (d *Database) Close() error {
	var err error
	if d.Pool != nil {
		d.Pool.Close()
	}
	if d.DB != nil {
		err = d.DB.Close()
	}
	return err
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Pool.Ping(ctx)
}

func (d *Database) GetPool() *pgxpool.Pool {
	return d.Pool
}

func (d *Database) GetDB() *sql.DB {
	return d.DB
}

func (d *Database) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := d.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}