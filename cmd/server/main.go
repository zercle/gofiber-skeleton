package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/database"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/server"

	_ "github.com/zercle/gofiber-skeleton/docs"
)

// @title Go Fiber Skeleton API
// @version 1.0
// @description Production-ready Go backend template using Fiber v2 framework with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize logger first
	logger.Init()

	// Create fx application
	app := fx.New(
		// Provide config
		fx.Provide(
			provideConfig,
			provideDatabase,
			server.NewFiberApp,
		),
		// Invoke server lifecycle
		fx.Invoke(runServer),
	)

	// Start the application
	app.Run()
}

// provideConfig loads and provides configuration
func provideConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

// provideDatabase connects to PostgreSQL and runs migrations
func provideDatabase(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
	log := logger.GetLogger()

	// Connect to database
	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connected successfully")

	// Run migrations automatically
	if err := database.RunMigrations(db, "db/migrations"); err != nil {
		log.Warn().Err(err).Msg("Failed to run migrations (this is optional during development)")
	} else {
		log.Info().Msg("Database migrations completed")
	}

	// Add lifecycle hooks
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Closing database connection...")
			return db.Close()
		},
	})

	return db, nil
}

// runServer starts the Fiber server with graceful shutdown
func runServer(lc fx.Lifecycle, app *server.FiberApp, cfg *config.Config) {
	log := logger.GetLogger()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Start server in a goroutine
			go func() {
				addr := fmt.Sprintf(":%s", cfg.Server.Port)
				log.Info().Msgf("Server starting on %s (env: %s)", addr, cfg.Server.Env)

				if err := app.App.Listen(addr); err != nil {
					log.Error().Err(err).Msg("Server failed to start")
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Shutting down server...")

			// Create a deadline for shutdown
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Shutdown the server gracefully
			if err := app.App.ShutdownWithContext(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("Server forced to shutdown")
				return err
			}

			log.Info().Msg("Server stopped gracefully")
			return nil
		},
	})

	// Handle OS signals for graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Info().Msg("Received shutdown signal")
	}()
}
