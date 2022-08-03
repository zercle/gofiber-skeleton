package datasources

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// New SQLite creates a new database connection backed by a given SQLite.
func NewSQLite(dbname string) (dbConn *gorm.DB, err error) {
	if len(dbname) == 0 {
		dbname = "test.db"
	}
	return gorm.Open(sqlite.Open(dbname), &gorm.Config{})
}
