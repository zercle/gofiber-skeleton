package datasources

import (
	"crypto/tls"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Resources struct {
	// client
	FastHttpClient *fasthttp.Client
	// parser pool
	JsonParserPool *fastjson.ParserPool
	// database
	MainDbConn *gorm.DB
	LogDbConn  *mongo.Database
	// Redis storage
	RedisStorage *redis.Storage
	// Session storage
	SessStore *session.Store
	// JWT
	JwtResources *JwtResources
}

func InitResources(httpclient *fasthttp.Client, jsonParserPool *fastjson.ParserPool, mainDbConn *gorm.DB, logDbConn *mongo.Database, redisStorage *redis.Storage, jwtResources *JwtResources) *Resources {
	return &Resources{
		FastHttpClient: httpclient,
		JsonParserPool: jsonParserPool,
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
