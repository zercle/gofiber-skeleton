package datasources

import (
	"gorm.io/gorm"
)

var (
	// MariaDB connection
	MariaDbConn *gorm.DB

	// Postgre SQL connection
	// PostgreSqlConn *gorm.DB

	// SQL Server Connection
	// SqlServerConn  *gorm.DB

	// sqLite Connection
	ConnSqlLite *gorm.DB

	// MongoDB Connection
	// MongoDbConn *mongo.Database

	// Redis storage
	// RedisStorage *redis.Storage

	// Session storage
	// SessStore *session.Store
)
