package infrastructure

import (
	"crypto"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/storage/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/internal/datasources"
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

type JwtResources struct {
	JwtVerifyKey     crypto.PublicKey
	JwtSignKey       crypto.PrivateKey
	JwtSigningMethod jwt.SigningMethod
	JwtKeyfunc       jwt.Keyfunc
	JwtParser        *jwt.Parser
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

func InitJwt(privateKeyPath, publicKeyPath string) (jwtResources JwtResources, err error) {
	jwtResources.JwtSignKey, jwtResources.JwtVerifyKey, jwtResources.JwtSigningMethod, err = datasources.InitJwtLocalKey(privateKeyPath, publicKeyPath)
	jwtResources.JwtKeyfunc = func(token *jwt.Token) (publicKey interface{}, err error) {
		if jwtResources.JwtVerifyKey == nil {
			err = fmt.Errorf("JWTVerifyKey not init yet")
		}
		return jwtResources.JwtVerifyKey, err
	}
	jwtResources.JwtParser = jwt.NewParser()
	return
}

func (r *JwtResources) IsJwtActive(tokenStr string) (token *jwt.Token, isActive bool, err error) {
	if r.JwtParser == nil {
		err = fmt.Errorf("JwtParser not init yet")
		return
	}
	claims := jwt.RegisteredClaims{}
	token, _, err = r.JwtParser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		return
	}

	isActive = claims.VerifyExpiresAt(time.Now(), true)

	return
}
