package main

import (
	"errors"
	"fmt"
	"gofiber-skeleton/internal/infra/config"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run ./cmd/migrator <up|down>")
	}

	_ = godotenv.Load()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	migrationsPath := "db/migrations"
	dbPath := cfg.DATABASE_URL

	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create database directory '%s': %v", dbDir, err)
		}
	}

	migrationDSN := fmt.Sprintf("sqlite://%s", dbPath)
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(sourceURL, migrationDSN)
	if err != nil {
		log.Fatalf("Migration setup failed for DSN '%s': %v", migrationDSN, err)
	}

	command := os.Args[1]
	var migrationErr error

	switch command {
	case "up":
		log.Println("Running migrations up...")
		migrationErr = m.Up()
	case "down":
		log.Println("Running migrations down...")
		migrationErr = m.Down()
	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down'.", command)
	}

	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		log.Fatalf("Migration failed: %v", migrationErr)
	}

	if migrationErr == nil {
		log.Printf("Migrations '%s' completed successfully.", command)
	} else {
		log.Println("No new migrations to apply.")
	}
}
