package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zercle/gofiber-skeleton/internal/config"
	appLogger "github.com/zercle/gofiber-skeleton/pkg/logger"
)

// SetupMiddleware configures and returns all middleware for the application
func SetupMiddleware(app *fiber.App, cfg *config.Config) {
	// Recovery middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			appLogger.Error("Panic recovered",
				"error", e,
				"method", c.Method(),
				"path", c.Path(),
				"ip", c.IP(),
			)
		},
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.CORS.AllowOrigins, ","),
		AllowMethods: strings.Join(cfg.CORS.AllowMethods, ","),
		AllowHeaders: strings.Join(cfg.CORS.AllowHeaders, ","),
		AllowCredentials: true,
		MaxAge: 86400, // 24 hours
	}))

	// Request logging middleware
	if cfg.IsDevelopment() {
		app.Use(fiberLogger.New(fiberLogger.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
			TimeFormat: time.RFC3339,
			TimeZone:   "UTC",
		}))
	} else {
		// Production logging with structured format
		app.Use(customLogger())
	}

	// Request ID middleware
	app.Use(requestID())
}

// customLogger provides structured logging for production
func customLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request completion
		duration := time.Since(start)

		appLogger.Info("HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"referer", c.Get("Referer"),
		)

		return err
	}
}

// requestID adds a unique request ID to each request
func requestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID exists in headers
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			// Generate a simple request ID
			reqID = fmt.Sprintf("%d-%d", time.Now().UnixNano(), len(c.Path()))
		}

		// Set request ID in headers and context
		c.Set("X-Request-ID", reqID)
		c.Locals("requestID", reqID)

		// Add request ID to logger context - simplified for Go 1.25
		// Note: slog doesn't have Context method in this version

		return c.Next()
	}
}

// HealthCheckMiddleware provides basic health check middleware
func HealthCheckMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set health check headers
		c.Set("X-Health-Check", "true")
		c.Set("X-Timestamp", time.Now().Format(time.RFC3339))

		return c.Next()
	}
}