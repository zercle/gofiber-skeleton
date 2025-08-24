package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config" // Import config package
	orderhandler "github.com/zercle/gofiber-skeleton/internal/order/handler"
	producthandler "github.com/zercle/gofiber-skeleton/internal/product/handler"
	userhandler "github.com/zercle/gofiber-skeleton/internal/user/handler"
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

	// Initialize SQLC queries

	// Handlers are now initialized within their respective RegisterRoutes functions if needed

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
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// JWT middleware
	app.Use(infrastructure.JWTMiddleware(cfg.JWT.Secret)) // Use JWT secret from config

	// API routes

	// Register domain-specific routes
	producthandler.RegisterRoutes(app)
	orderhandler.RegisterRoutes(app)
	userhandler.RegisterRoutes(app)

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
