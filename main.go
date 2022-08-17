package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/viper"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/configs"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"github.com/zercle/gofiber-skelton/internal/routes"
)

// PrdMode from GO_ENV
var (
	PrdMode    bool
	version    string
	build      string
	sessConfig session.Config
	logConfig  logger.Config
	runEnv     string
)

func main() {
	var err error
	// Running flag
	if len(os.Getenv("ENV")) != 0 {
		runEnv = os.Getenv("ENV")
	} else {
		flagEnv := flag.String("env", "dev", "A config file name without .env")
		flag.Parse()
		runEnv = *flagEnv
	}
	if err := configs.LoadConfig(runEnv); err != nil {
		log.Panicf("error while loading the env:\n %+v", err)
	}

	PrdMode = (viper.GetString("app.env") == "production")

	// Init datasources
	err = initDatasources()
	if err != nil {
		log.Panicf("error while init datasources:\n %+v", err)
	}

	// datasource resources
	resources := datasources.InitResources(datasources.MariaDbConn)

	// Init app
	log.Printf("init app")
	app := fiber.New(fiber.Config{
		ErrorHandler:   customErrorHandler,
		ReadTimeout:    60 * time.Second,
		ReadBufferSize: 8 * 1024,
		Prefork:        PrdMode,
		// speed up json with segmentio/encoding
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Post config by env
	err = configApp()
	if err != nil {
		log.Panicf("error while config app:\n %+v", err)
	}

	// Logger middleware for Fiber that logs HTTP request/response details.
	app.Use(logger.New(logConfig))

	// Recover middleware for Fiber that recovers from panics anywhere in the stack chain and handles the control to the centralized ErrorHandler.
	app.Use(recover.New(recover.Config{EnableStackTrace: !PrdMode}))

	// CORS middleware for Fiber that that can be used to enable Cross-Origin Resource Sharing with various options.
	app.Use(cors.New())

	// Init session
	// datasources.SessStore = session.New(sessConfig)

	// set apiV1 router
	routerResources := routes.InitRouterResources(resources)
	routerResources.SetupRoutes(app)

	// Log GO_ENV
	log.Printf("Runtime ENV: %s", viper.GetString("app.env"))
	log.Printf("Version: %s", version)
	log.Printf("Build: %s", build)

	// Listen from a different goroutine
	go func() {
		if err := app.ListenTLS(":"+viper.GetString("app.port.https"), viper.GetString("app.path.cert"), viper.GetString("app.path.priv")); err != nil {
			log.Panic(err)
		}
	}()

	// Create channel to signify a signal being sent
	quit := make(chan os.Signal, 1)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	<-quit
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	// if datasources.RedisStorage != nil {
	// 	datasources.RedisStorage.Close()
	// }
	fmt.Println("Successful shutdown.")
}

var customErrorHandler = func(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := http.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}

	responseData := helpers.ResponseForm{
		Success: false,
		Errors: []*helpers.ResposeError{
			{
				Code:    code,
				Message: err.Error(),
			},
		},
	}

	// Return statuscode with error message
	err = c.Status(code).JSON(responseData)
	if err != nil {
		// In case the JSON fails
		log.Printf("customErrorHandler: %+v", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}

func initDatasources() (err error) {

	// maxDBConn := 8
	// maxDBIdle := 4
	// Pre config by env
	// if !PrdMode {
	// 	maxDBConn = 2
	// 	maxDBIdle = 1
	// }

	// Init database connection
	datasources.ConnSqlLite, err = datasources.NewSQLite("book.db")

	// Init database connection
	// Create connection to database
	// log.Printf("connecting to %s/%s", viper.GetString("db.main.host"), viper.GetString("db.main.db_name"))
	// connMariaDB, err := datasources.MariadbConfig{
	// Username:     viper.GetString("db.main.username"),
	// Password:     viper.GetString("db.main.password"),
	// Host:         viper.GetString("db.main.host"),
	// Port:         viper.GetString("db.main.port"),
	// 	UnixSocket:   "",
	// 	MaxIdleConns: maxDBIdle,
	// 	MaxOpenConns: maxDBConn,
	// 	ParseTime:    true,
	// }.NewMariadbDB(viper.GetString("db.main.db_name"))
	// if err != nil {
	// 	log.Panicf("Error Connect to database:\n %+v", err)
	// }
	// log.Printf("connected to %s/%s", viper.GetString("db.main.host"), viper.GetString("db.main.db_name"))

	// datasources.MariaDbConn = connMariaDB

	// Create connection to redis
	// log.Printf("connecting to %s/%s", viper.GetString("db.redis.host"), viper.GetString("db.redis.db_name"))
	// redisPort, err := strconv.Atoi(viper.GetString("db.redis.port"))
	// if err != nil {
	// 	redisPort = 6379
	// 	err = nil
	// }
	// redisDB, err := strconv.Atoi(viper.GetString("db.redis.db_name"))
	// if err != nil {
	// 	redisDB = 0
	// 	err = nil
	// }
	// datasources.RedisStore = redis.New(redis.Config{
	// 	Host:     viper.GetString("db.redis.host"),
	// 	Port:     redisPort,
	// 	Username: viper.GetString("db.redis.username"),
	// 	Password: viper.GetString("db.redis.password"),
	// 	Database: redisDB,
	// })
	// log.Printf("connected to %s/%s", viper.GetString("db.redis.host"), viper.GetString("db.redis.db_name"))

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

func configApp() (err error) {
	if PrdMode {
		sessConfig = session.Config{
			Expiration:     8 * time.Hour,
			KeyLookup:      fmt.Sprintf("%s:%s", "cookie", viper.GetString("app.name")),
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: "Strict",
			// Storage:        datasources.RedisStorage,
			CookiePath: "/",
		}
		logFileWriter := &datasources.LogFileWriter{LogPath: "./log/gofiber-skelton", PrintConsole: true}
		logConfig = logger.Config{
			Format:     "[${blue}${time}${reset}] ${status} - ${ip},${ips} ${method} ${host} ${url}\tUserAgent:	${ua}\tReferer: ${referer}\tAuthorization: ${header:Authorization}\tBytesReceived: ${bytesReceived}\tBytesSent: ${bytesSent}\tError: ${red}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
			Output:     logFileWriter,
		}
	} else {
		sessConfig = session.ConfigDefault
		logConfig = logger.Config{
			Format:     "[${blue}${time}${reset}] ${status} - ${ip},${ips} ${method} ${host} ${url}\nUserAgent:\t${ua}\nReferer:\t${referer}\nAuthorization:\t${header:Authorization}\nBytesReceived:\t${bytesReceived}\nBytesSent:\t${bytesSent}\nError:\t${red}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
		}
	}
	return
}
