package main

import (
	"fmt"
	"os"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/pkg/logger"
	"github.com/zercle/gofiber-skeleton/pkg/server"
)

// Build information
var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(logger.LoggerConfig{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})

	logger.Info("Application starting",
		"version", Version,
		"build_time", BuildTime,
		"git_commit", GitCommit,
		"environment", cfg.Server.Environment,
	)

	// Create and start server
	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	logger.Info("Application shutdown completed")
}