package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func NewDatabase(databaseURL string) (*gorm.DB, error) {
	sqlDB, err := sql.Open("sqlite", databaseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to database successfully!")
	return db, nil
}