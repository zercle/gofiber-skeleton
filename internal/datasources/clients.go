package datasources

import (
	"crypto/tls"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
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

var (
	// Session storage
	SessStore *session.Store
	// parser pool
	JsonParserPool *fastjson.ParserPool
	// JWT
	JwtParser  *jwt.Parser
	JwtKeyfunc jwt.Keyfunc
)

func InitResources(fasthttpClient *fasthttp.Client, mainDbConn *gorm.DB, logDbConn *mongo.Database, redisStorage *redis.Storage, jwtResources *JwtResources) *Resources {
	return &Resources{
		FastHttpClient: fasthttpClient,
		MainDbConn:     mainDbConn,
		LogDbConn:      logDbConn,
		RedisStorage:   redisStorage,
		JwtResources:   jwtResources,
	}
}

func InitFastHttpClient(insecureSkipVerify bool) (client *fasthttp.Client) {
	client = &fasthttp.Client{
		MaxConnsPerHost: (runtime.NumCPU() * 512) / 2,
		ReadTimeout:     time.Second * 45,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}
	return
}

func InitJsonParserPool() (jsonParserPool *fastjson.ParserPool) {
	jsonParserPool = new(fastjson.ParserPool)
	return
}
