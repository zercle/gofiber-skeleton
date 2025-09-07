package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func main() {
	var direction = flag.String("direction", "up", "Migration direction: up, down, force, version")
	var steps = flag.Int("steps", -1, "Number of migration steps (default: all)")
	var version = flag.Uint("version", 0, "Force database to specific version")
	flag.Parse()

	cfg := config.NewConfig()

	m, err := migrate.New(
		"file://db/migrations",
		cfg.DatabaseURL(),
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		if *steps == -1 {
			err = m.Up()
			if err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run up migrations: %v", err)
			}
			fmt.Println("Migrations up completed successfully")
		} else {
			err = m.Steps(*steps)
			if err != nil {
				log.Fatalf("Failed to run %d migration steps: %v", *steps, err)
			}
			fmt.Printf("Migration %d steps completed successfully\n", *steps)
		}

	case "down":
		if *steps == -1 {
			err = m.Down()
			if err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run down migrations: %v", err)
			}
			fmt.Println("Migrations down completed successfully")
		} else {
			err = m.Steps(-*steps)
			if err != nil {
				log.Fatalf("Failed to run %d migration steps down: %v", *steps, err)
			}
			fmt.Printf("Migration %d steps down completed successfully\n", *steps)
		}

	case "force":
		if *version == 0 {
			fmt.Println("Version is required for force command")
			os.Exit(1)
		}
		err = m.Force(int(*version))
		if err != nil {
			log.Fatalf("Failed to force migration to version %d: %v", *version, err)
		}
		fmt.Printf("Migration forced to version %d\n", *version)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)

	default:
		fmt.Printf("Unknown direction: %s. Use: up, down, force, or version\n", *direction)
		os.Exit(1)
	}
}
