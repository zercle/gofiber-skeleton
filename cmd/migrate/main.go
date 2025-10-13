package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		action = flag.String("action", "", "Migration action: up, down, version")
		url    = flag.String("url", "", "Database URL (optional, defaults to MIGRATE_URL env var)")
		steps  = flag.Int("steps", 0, "Number of steps to migrate (for up/down)")
	)
	flag.Parse()

	if *action == "" {
		fmt.Println("Usage: migrate -action <up|down|version> [-url <database_url>] [-steps <n>]")
		fmt.Println("Environment variable MIGRATE_URL can be used instead of -url")
		os.Exit(1)
	}

	// Get database URL from parameter or environment
	dbURL := *url
	if dbURL == "" {
		dbURL = os.Getenv("MIGRATE_URL")
	}
	if dbURL == "" {
		// Fallback to default for development
		dbURL = "postgres://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable"
	}

	// Create migration instance
	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	defer func() {
		if err, _ := m.Close(); err != nil {
			log.Printf("Failed to close migration instance: %v", err)
		}
	}()

	// Execute migration action
	switch *action {
	case "up":
		if *steps > 0 {
			if err := m.Steps(*steps); err != nil {
				log.Fatalf("Failed to migrate up %d steps: %v", *steps, err)
			}
			fmt.Printf("Successfully migrated up %d steps\n", *steps)
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to migrate up: %v", err)
			}
			fmt.Println("Successfully migrated up")
		}

	case "down":
		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil {
				log.Fatalf("Failed to migrate down %d steps: %v", *steps, err)
			}
			fmt.Printf("Successfully migrated down %d steps\n", *steps)
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to migrate down: %v", err)
			}
			fmt.Println("Successfully migrated down")
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Println("No migrations applied")
			} else {
				log.Fatalf("Failed to get version: %v", err)
			}
		} else {
			status := "clean"
			if dirty {
				status = "dirty"
			}
			fmt.Printf("Current version: %d (%s)\n", version, status)
		}

	default:
		log.Fatalf("Unknown action: %s. Use: up, down, version", *action)
	}
}