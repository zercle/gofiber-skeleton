package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic recovered",
					slog.String("request_id", getRequestID(c)),
					slog.String("method", c.Method()),
					slog.String("path", c.Path()),
					slog.Any("panic", r),
				)

				// Return 500 error
				err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"code":    "INTERNAL_ERROR",
					"message": "Internal server error",
				})
			}
		}()

		return c.Next()
	}
}

// getRequestID retrieves request ID from context
func getRequestID(c *fiber.Ctx) string {
	if id := c.Locals("request_id"); id != nil {
		return fmt.Sprintf("%v", id)
	}
	return "unknown"
}
