package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/joho/godotenv"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"github.com/zercle/gofiber-skelton/internal/routes"
)

// prdMode from GO_ENV
var (
	prdMode    bool
	version    string
	build      string
	sessConfig session.Config
	logConfig  logger.Config
)

func main() {
	// Running flag
	runEnv := flag.String("env", "dev", "A env file name without .env")
	flag.Parse()
	// Load env
	err := godotenv.Load(*runEnv + ".env")
	if err != nil {
		log.Fatalf("error while loading the env:\n %+v", err)
	}

	prdMode = (os.Getenv("GO_ENV") == "production")

	maxDBConn := 8
	maxDBIdle := 4
	// Pre config by env
	if prdMode {
	} else {
		maxDBConn = 2
		maxDBIdle = 1
	}

	// Init database connection
	// Create connection to database
	datasources.MariaDB, err = datasources.MariadbConfig{
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		UnixSocket:   "",
		MaxIdleConns: maxDBIdle,
		MaxOpenConns: maxDBConn,
		ParseTime:    true,
	}.NewMariadbDB(os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("Error Connect to database:\n %+v", err)
	}

	// close the database connection if application errored.
	defer datasources.MariaDB.Close()

	// Create connection to redis
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		redisPort = 6379
		err = nil
	}
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
		err = nil
	}
	datasources.RedisStore = redis.New(redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: redisDB,
	})

	// close the redis connection if application errored.
	defer datasources.RedisStore.Close()

	// Init JWT Key
	datasources.JWTSignKey, datasources.JWTVerifyKey, datasources.JWTSigningMethod, err = datasources.JTWLocalKey(os.Getenv("JWT_PRIVATE"), os.Getenv("JWT_PUBLIC"))
	if err != nil {
		log.Fatalf("Error Init JWT Keys:\n %+v", err)
	}

	// Init Client
	datasources.FasthttpClient = datasources.InitFasthttpClient()
	datasources.HttpClient = datasources.InitHttpClient()
	datasources.RegxNum = regexp.MustCompile(`[0-9]+`)

	// Init app
	app := fiber.New(fiber.Config{
		ErrorHandler:   customErrorHandler,
		ReadTimeout:    60 * time.Second,
		ReadBufferSize: 8190,
		// Prefork:      prdMode,
	})

	// Post config by env
	if prdMode {
		sessConfig = session.Config{
			Expiration:     8 * time.Hour,
			KeyLookup:      fmt.Sprintf("%s:%s", "cookie", os.Getenv("SESS_ID")),
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: "Strict ",
			Storage:        datasources.RedisStore,
			CookiePath:     "/",
		}
		logFileWriter := &datasources.LogFileWriter{LogPath: "/var/log/gofiber-skelton", PrintConsole: true}
		logConfig = logger.Config{
			Format: "[${blue}${time}${reset}] ${status} - ${ip},${ips} ${method} ${host} ${url}\tUserAgent:	${ua}\tReferer: ${referer}\tAuthorization: ${header:Authorization}\tBytesReceived: ${bytesReceived}\tBytesSent: ${bytesSent}\tError: ${red}${error}${reset}\n",
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

	// Logger middleware for Fiber that logs HTTP request/response details.
	app.Use(logger.New(logConfig))

	// Recover middleware for Fiber that recovers from panics anywhere in the stack chain and handles the control to the centralized ErrorHandler.
	app.Use(recover.New())

	// CORS middleware for Fiber that that can be used to enable Cross-Origin Resource Sharing with various options.
	app.Use(cors.New())

	// Init session
	datasources.SessStore = session.New(sessConfig)

	// set api router
	routersV1DBHandler := routes.InitRouterResources(datasources.MariaDB)
	routersV1DBHandler.SetRouters(app)

	// Log GO_ENV
	log.Printf("Runtime ENV: %s", os.Getenv("GO_ENV"))
	log.Printf("Version: %s", version)
	log.Printf("Build: %s", build)

	// Listen from a different goroutine
	// go func() {
	// 	if err := app.Listen(os.Getenv("HTTP_PORT")); err != nil {
	// 		log.Panic(err)
	// 	}
	// }()
	go func() {
		if err := app.ListenTLS(os.Getenv("HTTPS_PORT"), os.Getenv("CERT_PATH"), os.Getenv("PRIV_PATH")); err != nil {
			log.Panic(err)
		}
	}()

	// Create channel to signify a signal being sent
	ch := make(chan os.Signal, 2)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	<-ch
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
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
