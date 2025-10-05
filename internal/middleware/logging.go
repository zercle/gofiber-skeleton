package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/logger"
)

// StructuredLogger logs HTTP requests with structured fields
func StructuredLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Get request ID from context
		requestID, ok := c.Locals("request_id").(string)
		if !ok {
			requestID = ""
		}

		// Log structured request
		log := logger.GetLogger().Info().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Int("body_size", len(c.Response().Body()))

		// Add error if present
		if err != nil {
			log = log.Err(err)
		}

		// Add user agent
		if ua := c.Get("User-Agent"); ua != "" {
			log = log.Str("user_agent", ua)
		}

		log.Msg("HTTP request")

		return err
	}
}
