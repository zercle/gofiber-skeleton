package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"

	"github.com/zercle/gofiber-skeleton/internal/cache"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/internal/response"

	userHandler "github.com/zercle/gofiber-skeleton/internal/user/handler"
	userRepository "github.com/zercle/gofiber-skeleton/internal/user/repository"
	userUsecase "github.com/zercle/gofiber-skeleton/internal/user/usecase"

	postHandler "github.com/zercle/gofiber-skeleton/internal/post/handler"
	postRepository "github.com/zercle/gofiber-skeleton/internal/post/repository"
	postUsecase "github.com/zercle/gofiber-skeleton/internal/post/usecase"

	fiberswagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/zercle/gofiber-skeleton/docs" // Required for swagger documentation generation
)

// FiberApp wraps the Fiber application
type FiberApp struct {
	App *fiber.App
}

// NewFiberApp creates and configures the Fiber application using fx
func NewFiberApp(sqlDB *sql.DB, cfg *config.Config) *FiberApp {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			logger.GetLogger().Error().Err(err).Msg("Unhandled application error")
			return c.Status(code).SendString(err.Error())
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.StructuredLogger())
	app.Use(middleware.CORS())
	app.Use(middleware.Security())

	// Health check routes
	app.Get("/health", healthHandler)
	app.Get("/ready", readinessHandler(sqlDB))

	// Swagger documentation
	app.Get("/swagger/*", fiberswagger.New(fiberswagger.Config{
		URL:          "http://localhost:8080/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// API route group with rate limiting
	api := app.Group("/api", middleware.APIRateLimit())

	// Initialize database queries
	queries := db.New(sqlDB)

	// Optional: Initialize Redis client (gracefully handle if Redis is unavailable)
	var redisClient *cache.RedisClient
	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		logger.GetLogger().Warn().Err(err).Msg("Redis connection failed, continuing without cache")
		redisClient = nil
	} else {
		logger.GetLogger().Info().Msg("Redis connected successfully")
	}

	// Initialize user domain
	userRepo := userRepository.NewPostgresUserRepository(queries)
	authUsecase := userUsecase.NewAuthUsecase(userRepo, cfg.JWT.Secret)

	// Initialize post domain
	postRepo := postRepository.NewPostgresPostRepository(queries)
	postUsecaseInstance := postUsecase.NewPostUsecase(postRepo)

	// API v1 routes
	v1 := api.Group("/v1")

	// Auth routes with stricter rate limiting
	authRoutes := v1.Group("/auth", middleware.AuthRateLimit())
	userHandler.RegisterAuthRoutes(authRoutes, authUsecase)

	// Post routes
	postHandler.RegisterPostRoutes(v1, postUsecaseInstance, cfg.JWT.Secret)

	// Stub routes for future implementation
	v1.Get("/users", stubHandler)
	v1.Get("/threads", stubHandler)
	v1.Get("/comments", stubHandler)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return response.Fail(c, http.StatusNotFound, fiber.Map{
			"error": "Route not found",
		})
	})

	// Store Redis client in app context if available
	if redisClient != nil {
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("redis", redisClient)
			return c.Next()
		})
	}

	return &FiberApp{App: app}
}

// Module provides the Fiber app as an fx module
var Module = fx.Options(
	fx.Provide(NewFiberApp),
)

// stubHandler is a placeholder for unimplemented routes
func stubHandler(c *fiber.Ctx) error {
	logger.GetLogger().Warn().Msgf("Route %s not implemented", c.OriginalURL())
	return response.Fail(c, fiber.StatusNotImplemented, fiber.Map{
		"error":   "Not implemented",
		"message": fmt.Sprintf("Route %s is not yet implemented", c.OriginalURL()),
	})
}

// healthHandler checks if the service is alive
func healthHandler(c *fiber.Ctx) error {
	return response.Success(c, http.StatusOK, fiber.Map{
		"status": "healthy",
	})
}

// readinessHandler checks if the service is ready to accept requests
func readinessHandler(sqlDB *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check database connection
		if err := sqlDB.Ping(); err != nil {
			logger.GetLogger().Error().Err(err).Msg("Database health check failed")
			return response.Error(c, http.StatusServiceUnavailable, "Service not ready", 5001)
		}

		// Check Redis if available
		redisStatus := "not configured"
		if redisValue := c.Locals("redis"); redisValue != nil {
			if redis, ok := redisValue.(*cache.RedisClient); ok && redis != nil {
				if err := redis.GetClient().Ping(c.Context()).Err(); err != nil {
					redisStatus = "disconnected"
				} else {
					redisStatus = "connected"
				}
			}
		}

		return response.Success(c, http.StatusOK, fiber.Map{
			"status":   "ready",
			"database": "connected",
			"redis":    redisStatus,
		})
	}
}
