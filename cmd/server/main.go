package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/samber/do/v2"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/container"
	postdelivery "github.com/zercle/gofiber-skeleton/internal/domains/post/delivery"
	userdelivery "github.com/zercle/gofiber-skeleton/internal/domains/user/delivery"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/utils"
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

	// Setup dependency injection container
	diContainer := container.NewContainer(cfg)

	// Create Fiber app
	app := fiber.New()

	// Setup middleware
	setupMiddleware(app, cfg)

	// Setup routes with DI
	setupRoutes(app, cfg, diContainer)

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

func setupRoutes(app *fiber.App, cfg *config.Config, diContainer do.Injector) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": c.Locals("request_id"),
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Get handlers from DI container
	userHandler := do.MustInvoke[*userdelivery.UserHandler](diContainer)
	postHandler := do.MustInvoke[*postdelivery.PostHandler](diContainer)

	// Register user routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// User routes (protected)
	users := api.Group("/users")
	users.Get("/profile", utils.AuthMiddleware(), userHandler.GetProfile)
	users.Put("/profile", utils.AuthMiddleware(), userHandler.UpdateProfile)
	users.Post("/change-password", utils.AuthMiddleware(), userHandler.ChangePassword)
	users.Get("/", utils.AuthMiddleware(), userHandler.ListUsers) // Admin only
	users.Delete("/:id", utils.AuthMiddleware(), userHandler.DeactivateUser) // Admin only

	// Register post routes
	postHandler.RegisterRoutes(app)

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
