package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func main() {
	var direction = flag.String("direction", "up", "Migration direction: up or down")
	var steps = flag.Int("steps", 0, "Number of migration steps (0 for all)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.Database.URL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Printf("Migration completed successfully")
	} else {
		log.Printf("Migration completed successfully. Current version: %d, dirty: %t", version, dirty)
	}
}