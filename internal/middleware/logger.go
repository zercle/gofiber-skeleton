package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LoggerMiddleware logs incoming requests with structured logging
func LoggerMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)

		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)
		log.Info("HTTP Request",
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration_ms", duration),
			slog.String("ip", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}
