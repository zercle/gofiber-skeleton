package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

// @title Go Fiber Skeleton API
// @version 1.0
// @description A production-ready Go Fiber template with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Note: Container setup would go here in full implementation

	// Create Fiber app
	app := fiber.New()

	// Setup middleware
	setupMiddleware(app, cfg)

	// Setup routes
	setupRoutes(app, cfg)

	// Start server
	go func() {
		addr := cfg.AppHost + ":" + cfg.AppPort
		log.Printf("Server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(app)
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Recovery middleware
	app.Use(middleware.Recovery())

	// Request ID middleware
	app.Use(middleware.RequestID())

	// Logger middleware
	app.Use(middleware.CustomLogger())

	// CORS middleware
	if cfg.IsDevelopment() {
		app.Use(middleware.DevelopmentCORS())
	} else {
		app.Use(middleware.RestrictedCORS([]string{"https://yourdomain.com"}))
	}

	// Error handling middleware
	app.Use(middleware.ErrorHandler())
}

func setupRoutes(app *fiber.App, cfg *config.Config) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": c.Locals("request_id"),
		})
	})

	// Note: For now, we'll create handlers without DI
	// In a real implementation, you would set up proper DI
	// userHandler := delivery.NewUserHandler(/* dependencies */)

	// API routes
	api := app.Group("/api/v1")

	// Simple test endpoint for now
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API is working!",
			"config": fiber.Map{
				"app_name": cfg.AppName,
				"app_env":  cfg.AppEnv,
			},
		})
	})

	// Setup 404 handler
	app.Use(middleware.NotFoundHandler())

	// Setup 405 handler
	app.Use(middleware.MethodNotAllowedHandler())
}

func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}