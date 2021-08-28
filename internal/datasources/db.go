package datasources

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

var (
	MariaDB    *MariadbConfig
	PostgreSQL *PostgreSQLConfig
	// RedisStore for session store
	RedisStore *redis.Storage
	// SessStore obj
	SessStore *session.Store
)
