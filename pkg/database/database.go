package database

import (
	"fmt"
	"log"

	

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(databaseURL string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully!")
	return db
}

func RunMigrations(migrationsPath string, databaseURL string) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("sqlite://%s", databaseURL),
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}