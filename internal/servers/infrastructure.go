package servers

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/gofiber/storage/redis"
	"github.com/spf13/viper"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"gorm.io/gorm"
)

func initDatasources() (err error) {
	// Init JWT Key
	log.Printf("init JWT")
	jwtSignKey, jwtVerifyKey, jwtSigningMethod, err := datasources.JTWLocalKey(viper.GetString("jwt.private"), viper.GetString("jwt.public"))
	if err != nil {
		log.Panicf("Error Init JWT Keys:\n %+v", err)
	}
	log.Printf("init JWT %s done", jwtSigningMethod.Alg())

	datasources.JWTSignKey = &jwtSignKey
	datasources.JWTVerifyKey = &jwtVerifyKey
	datasources.JWTSigningMethod = &jwtSigningMethod

	// Init Client
	log.Printf("init client")
	datasources.FastHttpClient = datasources.InitFasthttpClient()
	datasources.JsonParserPool = datasources.InitJsonParserPool()
	datasources.RegxNum = regexp.MustCompile(`[0-9]+`)
	log.Printf("init client done")

	return
}

func connectToSqLite() (dbConn *gorm.DB, err error) {
	return datasources.NewSQLite(viper.GetString("db.sqlite.db_name"))
}

func connectToMariaDb() (dbConn *gorm.DB, err error) {
	// Init database connection
	// Create connection to database
	log.Printf("connecting to %s/%s", viper.GetString("db.mariadb.host"), viper.GetString("db.mariadb.db_name"))
	dbConfig := datasources.MariadbConfig{
		Username:     viper.GetString("db.mariadb.username"),
		Password:     viper.GetString("db.mariadb.password"),
		Host:         viper.GetString("db.mariadb.host"),
		Port:         viper.GetString("db.mariadb.port"),
		MaxIdleConns: viper.GetInt("db.mariadb.conn.min"),
		MaxOpenConns: viper.GetInt("db.mariadb.conn.max"),
		ParseTime:    true,
	}
	if _, err := os.Stat(viper.GetString("db.mariadb.sock")); err == nil {
		dbConfig.UnixSocket = viper.GetString("db.mariadb.sock")
		log.Printf("connecting by %s", dbConfig.UnixSocket)
	}
	dbConn, err = dbConfig.NewMariaDB(viper.GetString("db.mariadb.db_name"))
	if err != nil {
		log.Panicf("Error Connect to database:\n %+v", err)
	}
	log.Printf("connected to %s/%s", viper.GetString("db.mariadb.host"), viper.GetString("db.mariadb.db_name"))

	return
}

func connectToPostgres() (dbConn *gorm.DB, err error) {
	// Init database connection
	// Create connection to database
	log.Printf("connecting to %s/%s", viper.GetString("db.postgres.host"), viper.GetString("db.postgres.db_name"))
	dbConfig := datasources.PostgreSQLConfig{
		Username:     viper.GetString("db.postgres.username"),
		Password:     viper.GetString("db.postgres.password"),
		Host:         viper.GetString("db.postgres.host"),
		Port:         viper.GetString("db.postgres.port"),
		MaxIdleConns: viper.GetInt("db.postgres.conn.min"),
		MaxOpenConns: viper.GetInt("db.postgres.conn.max"),
		ParseTime:    true,
	}
	if _, err := os.Stat(viper.GetString("db.postgres.sock")); err == nil {
		dbConfig.UnixSocket = viper.GetString("db.postgres.sock")
		log.Printf("connecting by %s", dbConfig.UnixSocket)
	}
	dbConn, err = dbConfig.NewPostgreSQL(viper.GetString("db.postgres.db_name"))
	if err != nil {
		log.Panicf("Error Connect to database:\n %+v", err)
	}
	log.Printf("connected to %s/%s", viper.GetString("db.postgres.host"), viper.GetString("db.postgres.db_name"))

	return
}

func connectToRedis() (dbConn *redis.Storage, err error) {
	// Create connection to redis
	log.Printf("connecting to %s/%s", viper.GetString("db.redis.host"), viper.GetString("db.redis.db_name"))
	redisPort, err := strconv.Atoi(viper.GetString("db.redis.port"))
	if err != nil {
		redisPort = 6379
		err = nil
	}
	redisDB, err := strconv.Atoi(viper.GetString("db.redis.db_name"))
	if err != nil {
		redisDB = 0
		err = nil
	}
	dbConn = redis.New(redis.Config{
		Host:     viper.GetString("db.redis.host"),
		Port:     redisPort,
		Username: viper.GetString("db.redis.username"),
		Password: viper.GetString("db.redis.password"),
		Database: redisDB,
	})
	log.Printf("connected to %s/%s", viper.GetString("db.redis.host"), viper.GetString("db.redis.db_name"))
	return
}
