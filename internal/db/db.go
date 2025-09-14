package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB initializes and returns a GORM DB instance
func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}