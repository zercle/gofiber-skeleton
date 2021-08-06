package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"github.com/zercle/gofiber-skelton/internal/routes"
)

// PrdMode from GO_ENV
var (
	PrdMode    bool
	version    string
	build      string
	sessConfig session.Config
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

	PrdMode = (os.Getenv("GO_ENV") == "production")

	// Pre config by env
	if PrdMode {
	} else {
	}

	// Init Client
	datasources.FasthttpClient = datasources.InitFasthttpClient()
	datasources.HttpClient = datasources.InitHttpClient()

	// Init app
	app := fiber.New(fiber.Config{
		ErrorHandler:   customErrorHandler,
		ReadTimeout:    60 * time.Second,
		ReadBufferSize: 8190,
		// Prefork:      PrdMode,
	})

	// Post config by env
	if PrdMode {
		sessConfig = session.Config{
			Expiration:     8 * time.Hour,
			KeyLookup:      fmt.Sprintf("%s:%s", "cookie", os.Getenv("SESS_ID")),
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: "Strict ",
			Storage:        datasources.RedisStore,
			CookiePath:     "/",
		}
	} else {
		sessConfig = session.ConfigDefault
		app.Use(logger.New())
	}

	// Recover middleware for Fiber that recovers from panics anywhere in the stack chain and handles the control to the centralized ErrorHandler.
	app.Use(recover.New())

	// CORS middleware for Fiber that that can be used to enable Cross-Origin Resource Sharing with various options.
	app.Use(cors.New())

	// Init session
	datasources.SessStore = session.New(sessConfig)

	// set api router
	routersV1DBHandler := routes.InitRouterResources(datasources.MainDB)
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
