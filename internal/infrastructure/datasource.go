package infrastructure

import (
	"github.com/spf13/viper"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"gorm.io/gorm"
)

func connectToSqLite() (dbConn *gorm.DB, err error) {
	return datasources.NewSQLite(viper.GetString("db.sqlite.db_name"))
}
