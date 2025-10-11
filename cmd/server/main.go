package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/swagger"
	"github.com/zercle/gofiber-skeleton/internal/config"
	userDelivery "github.com/zercle/gofiber-skeleton/internal/domains/user/delivery"
	userRepo "github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	userUsecase "github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/auth"
	"github.com/zercle/gofiber-skeleton/pkg/cache"
	"github.com/zercle/gofiber-skeleton/pkg/database"
	"github.com/zercle/gofiber-skeleton/pkg/response"
	"github.com/zercle/gofiber-skeleton/pkg/validator"

	_ "github.com/zercle/gofiber-skeleton/docs" // Import generated docs
)

// @title Go Fiber Skeleton API
// @version 1.0
// @description Production-ready Go Fiber backend template with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

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
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Database connected successfully")

	// Initialize cache
	cacheClient, err := cache.NewValkey(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to cache: %v", err)
		log.Println("Continuing without cache...")
	} else {
		defer cacheClient.Close()
		log.Println("✓ Cache connected successfully")
	}

	// Initialize utilities
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	logger := middleware.NewCustomLogger(cfg)

	// Initialize repositories
	userRepository := userRepo.NewPostgresUserRepository(db.Pool)

	// Initialize use cases
	userUsecaseInstance := userUsecase.NewUserUsecase(
		userRepository,
		jwtManager,
		validatorInstance,
		cfg,
	)

	// Initialize handlers
	userHandler := userDelivery.NewUserHandler(userUsecaseInstance)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Go Fiber Skeleton",
		ServerHeader:          "Go Fiber Skeleton",
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		ErrorHandler:          errorHandler,
		DisableStartupMessage: false,
	})

	// Setup middleware
	app.Use(middleware.Recover())
	app.Use(middleware.Security())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(cfg))
	app.Use(middleware.CORS(cfg))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Apply rate limiting to API routes
	apiGroup := app.Group("/api", middleware.RateLimit(cfg))

	// Health check endpoint (no rate limit)
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "healthy",
			"service": "go-fiber-skeleton",
			"time":    time.Now().UTC(),
		}

		// Check database health
		if err := db.Health(c.Context()); err != nil {
			health["database"] = "unhealthy"
			health["status"] = "degraded"
		} else {
			health["database"] = "healthy"
		}

		// Check cache health
		if cacheClient != nil {
			if err := cacheClient.Health(c.Context()); err != nil {
				health["cache"] = "unhealthy"
			} else {
				health["cache"] = "healthy"
			}
		}

		return c.JSON(health)
	})

	// Swagger documentation
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
		log.Println("✓ Swagger documentation available at /swagger/")
	}

	// API v1 routes
	v1 := apiGroup.Group("/v1")

	// Register domain routes
	userDelivery.RegisterRoutes(v1, userHandler, cfg)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c)
	})

	// Start server
	logger.Info("Starting server", map[string]interface{}{
		"address": cfg.GetServerAddress(),
		"env":     cfg.Server.Env,
	})

	// Graceful shutdown
	go func() {
		if err := app.Listen(cfg.GetServerAddress()); err != nil {
			logger.Error("Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logger.Info("Server exited")
}

// errorHandler handles errors globally
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(response.Response{
		Status: "error",
		Error: &response.Error{
			Code:    fmt.Sprintf("ERROR_%d", code),
			Message: message,
		},
	})
}
