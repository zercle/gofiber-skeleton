package datasources

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	ConnMariaDb *gorm.DB
	// ConnPostgreSql *gorm.DB
	// ConnSqlServer  *gorm.DB
	// ConnSqlLite    *gorm.DB
	ConnMongoDb *mongo.Database
	// RedisStore for session store
	RedisStore *redis.Storage
	// SessStore obj
	SessStore *session.Store
)
