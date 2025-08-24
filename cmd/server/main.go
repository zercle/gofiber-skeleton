package main

import (
	"log"
	"time"

	// @title E-commerce API
	// @version 1.0
	// @description This is a sample API for an e-commerce application.

	// @contact.name API Support
	// @contact.url https://zercle.tech
	// @contact.email system-admin@zercle.tech

	// @license.name MIT
	// @license.url https://mit-license.org/

	// @host localhost:8080
	// @BasePath /api/v1
	// @schemes http
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/app"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config" // Import config package
	orderhandler "github.com/zercle/gofiber-skeleton/internal/order/handler"
	producthandler "github.com/zercle/gofiber-skeleton/internal/product/handler"
	userhandler "github.com/zercle/gofiber-skeleton/internal/user/handler"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/samber/do/v2"
	_ "github.com/zercle/gofiber-skeleton/docs" // Import generated docs
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := infrastructure.ConnectDatabase(cfg.DB) // Pass database URL from config
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new Injector
	injector := app.NewInjector(db, &cfg) // Pass pointer to cfg
	defer func(i *do.RootScope) {
		if i != nil {
			if errMap := i.Shutdown().Errors; len(errMap) != 0 {
				log.Printf("injector shutdown error: %v", err)
			}
		}
	}(injector)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		ReadTimeout: 1 * time.Minute,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Swagger API docs
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:  false,
		DocExpansion: "none",
		URL:          "/swagger/doc.json", // Correct URL to the generated swagger.json
	}))

	// API routes
	apiV1Group := app.Group("/api/v1")

	// Register domain-specific routes
	// Resolve ProductHandler and initialize routes
	productHandler := do.MustInvoke[*producthandler.ProductHandler](injector)

	// Resolve OrderHandler and initialize routes
	orderHandler := do.MustInvoke[*orderhandler.OrderHandler](injector)

	// Resolve UserHandler and initialize routes
	userHandler := do.MustInvoke[*userhandler.UserHandler](injector)

	// Protected routes group
	jwtHandler := infrastructure.JWTMiddleware(cfg.JWT.Secret)

	userhandler.SetupRoutes(apiV1Group, jwtHandler, userHandler)
	producthandler.SetupRoutes(apiV1Group, jwtHandler, productHandler)
	orderhandler.SetupRoutes(apiV1Group, jwtHandler, orderHandler)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "E-commerce API is running",
		})
	})

	// Start server in other goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.App.Port)
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quitCh := make(chan os.Signal, 1)

	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	<-quitCh

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	// Clean up connections

	log.Println("Server gracefully shutdown")

}
