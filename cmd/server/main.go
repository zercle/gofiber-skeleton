package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/server"

	_ "github.com/zercle/gofiber-skeleton/docs"
)

func main() {
	// Initialize logger
	logger.Init()
	log := logger.GetLogger()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to PostgreSQL database using GORM
	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get underlying sql.DB from gorm.DB")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	// Initialize Fiber app and register routes
	app := server.SetupRouter(gormDB, cfg)

	// Start server
	log.Fatal().Err(app.Listen(fmt.Sprintf(":%s", cfg.Port))).Msg("Server stopped")
}