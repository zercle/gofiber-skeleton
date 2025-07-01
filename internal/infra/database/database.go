package database

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func ConnectDB(databaseURL string) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", databaseURL)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully!")
	return db
}