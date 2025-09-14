package main

import (
	"fmt"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/server"
)

func main() {
	// Initialize logger
	logger.Init()
	log := logger.GetLogger()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.ConnectDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		sqlDB, err := database.DB()
		if err != nil {
			log.Error().Err(err).Msg("Error getting sql.DB from gorm.DB")
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	// Initialize Fiber app and register routes
	app := server.SetupRouter(database, cfg)

	// Start server
	log.Fatal().Err(app.Listen(fmt.Sprintf(":%s", cfg.Port))).Msg("Server stopped")
}