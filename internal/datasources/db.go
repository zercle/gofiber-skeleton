package datasources

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"gorm.io/gorm"
)

var (
	MariaDB    *gorm.DB
	PostgreSQL *gorm.DB
	// RedisStore for session store
	RedisStore *redis.Storage
	// SessStore obj
	SessStore *session.Store
)
