package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimit returns rate limiting middleware
func RateLimit(max int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: duration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// AuthRateLimit returns rate limiting for authentication endpoints
func AuthRateLimit() fiber.Handler {
	return RateLimit(5, 15*time.Minute) // 5 requests per 15 minutes
}

// APIRateLimit returns rate limiting for API endpoints
func APIRateLimit() fiber.Handler {
	return RateLimit(100, 1*time.Minute) // 100 requests per minute
}
