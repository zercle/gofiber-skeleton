package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/zercle/gofiber-skeleton/internal/response"
)

// RateLimit creates a rate limiting middleware
func RateLimit(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return response.Fail(c, fiber.StatusTooManyRequests, fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
	})
}

// APIRateLimit creates a standard rate limit for API endpoints
func APIRateLimit() fiber.Handler {
	// 100 requests per minute
	return RateLimit(100, 1*time.Minute)
}

// AuthRateLimit creates a stricter rate limit for authentication endpoints
func AuthRateLimit() fiber.Handler {
	// 5 requests per minute to prevent brute force attacks
	return RateLimit(5, 1*time.Minute)
}