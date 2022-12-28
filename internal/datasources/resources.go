package datasources

import (
	"github.com/gofiber/storage/redis"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Resources struct {
	// client
	FastHttpClient *fasthttp.Client
	// database
	MainDbConn *gorm.DB
	LogDbConn  *mongo.Database
	// Redis storage
	RedisStorage *redis.Storage
	// JWT
	JwtResources *JwtResources
}

func InitResources(fasthttpClient *fasthttp.Client, mainDbConn *gorm.DB, logDbConn *mongo.Database, redisStorage *redis.Storage, jwtResources *JwtResources) *Resources {
	return &Resources{
		FastHttpClient: fasthttpClient,
		MainDbConn:     mainDbConn,
		LogDbConn:      logDbConn,
		RedisStorage:   redisStorage,
		JwtResources:   jwtResources,
	}
}
