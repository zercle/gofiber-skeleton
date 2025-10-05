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
			// Generate new UUIDv7 if not present
			id, err := uuid.NewV7()
			if err != nil {
				// Fallback to empty string if UUID generation fails (extremely rare)
				requestID = ""
			} else {
				requestID = id.String()
			}
		}

		// Set request ID in context for use in handlers
		c.Locals("request_id", requestID)

		// Set request ID in response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
