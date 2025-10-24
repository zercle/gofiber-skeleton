package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/do/v2"
	db "github.com/zercle/template-go-fiber/internal/infrastructure/sqlc"
	"github.com/zercle/template-go-fiber/internal/repositories"
)

// InitializeDI initializes the dependency injection container
func InitializeDI(cfg *Config) (do.Injector, error) {
	injector := do.New()

	// Register config
	do.ProvideValue(injector, cfg)

	// Register database connection
	do.Provide(injector, func(i do.Injector) (*sql.DB, error) {
		cfg := do.MustInvoke[*Config](i)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)

		database, err := sql.Open(cfg.Database.Driver, dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		// Configure connection pool
		database.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		database.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		database.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

		// Test the connection
		if err := database.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}

		return database, nil
	})

	// Register sqlc Queries
	do.Provide(injector, func(i do.Injector) (*db.Queries, error) {
		database := do.MustInvoke[*sql.DB](i)
		return db.New(database), nil
	})

	// Register repositories
	do.Provide(injector, func(i do.Injector) (*repositories.UserRepository, error) {
		queries := do.MustInvoke[*db.Queries](i)
		return repositories.NewUserRepository(queries), nil
	})

	return injector, nil
}
