package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists in header
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			// Generate new UUID if not present
			requestID = uuid.New().String()
		}

		// Set request ID in context for use in handlers
		c.Locals("request_id", requestID)

		// Set request ID in response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
