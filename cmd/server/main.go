package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/database"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/handler"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

// @title Go Fiber Skeleton API
// @version 1.0
// @description Production-ready Go Fiber template with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db.DB, "db/migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ServerHeader: "Fiber",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: customErrorHandler,
	})

	// Setup middleware
	app.Use(middleware.Recovery())
	app.Use(middleware.RequestID())
	app.Use(middleware.CORS())
	app.Use(middleware.Security())
	if cfg.App.IsDevelopment() {
		app.Use(middleware.DetailedLogger())
	} else {
		app.Use(middleware.ProductionLogger())
	}

	// Health check endpoints
	app.Get("/health", healthCheck(db))
	app.Get("/health/ready", readinessCheck(db))
	app.Get("/health/live", livenessCheck())

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db.DB)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)

	// API routes
	api := app.Group("/api/v1")
	api.Use(middleware.APIRateLimit())

	// Setup domain routes
	handler.SetupUserRoutes(api, authUsecase)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s (environment: %s)", addr, cfg.App.Environment)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"code":    code,
	})
}

// healthCheck returns overall health status
func healthCheck(db *database.PostgresDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check database
		if err := db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "down",
				"error":    err.Error(),
			})
		}

		stats := db.GetStats()
		return c.JSON(fiber.Map{
			"status": "healthy",
			"database": fiber.Map{
				"status":          "up",
				"open_connections": stats.OpenConnections,
				"in_use":          stats.InUse,
				"idle":            stats.Idle,
			},
		})
	}
}

// readinessCheck returns readiness status for Kubernetes
func readinessCheck(db *database.PostgresDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"ready": false,
			})
		}
		return c.JSON(fiber.Map{
			"ready": true,
		})
	}
}

// livenessCheck returns liveness status for Kubernetes
func livenessCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"alive": true,
		})
	}
}
