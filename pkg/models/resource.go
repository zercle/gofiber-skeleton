package models

import (
	"crypto"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Resources struct {
	SessConfig session.Config
	LogConfig  logger.Config
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

type JwtResources struct {
	JwtVerifyKey     crypto.PublicKey
	JwtSignKey       crypto.PrivateKey
	JwtSigningMethod jwt.SigningMethod
	JwtKeyfunc       jwt.Keyfunc
	JwtParser        *jwt.Parser
}
