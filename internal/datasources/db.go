package datasources

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	ConnMariaDB    *MariaDB
	ConnPostgreSQL *PostgreSQL
	ConnMongoDB    *MongoDB
	// RedisStore for session store
	RedisStore *redis.Storage
	// SessStore obj
	SessStore *session.Store
)

type MariaDB struct {
	*gorm.DB
}

type PostgreSQL struct {
	*gorm.DB
}

type MongoDB struct {
	*mongo.Database
}
