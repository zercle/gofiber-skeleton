package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
)

// DistributedRateLimiter creates a Redis-backed distributed rate limiter
func DistributedRateLimiter(redisHost string, redisPort int, max int, duration time.Duration) fiber.Handler {
	storage := redis.New(redis.Config{
		Host:     redisHost,
		Port:     redisPort,
		Database: 0,
		Reset:    false,
	})

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: duration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Rate limit exceeded. Please try again later.",
				"retry_after": int(duration.Seconds()),
			})
		},
		Storage: storage,
	})
}

// DistributedAuthRateLimit returns distributed rate limiting for auth endpoints
func DistributedAuthRateLimit(redisHost string, redisPort int) fiber.Handler {
	return DistributedRateLimiter(redisHost, redisPort, 5, 15*time.Minute)
}

// DistributedAPIRateLimit returns distributed rate limiting for API endpoints
func DistributedAPIRateLimit(redisHost string, redisPort int) fiber.Handler {
	return DistributedRateLimiter(redisHost, redisPort, 100, 1*time.Minute)
}

// PerUserRateLimiter creates a rate limiter based on user ID
func PerUserRateLimiter(redisHost string, redisPort int, max int, duration time.Duration) fiber.Handler {
	storage := redis.New(redis.Config{
		Host:     redisHost,
		Port:     redisPort,
		Database: 0,
	})

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: duration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Extract user ID from JWT claims or context
			userID := c.Locals("userID")
			if userID == nil {
				return c.IP() // Fallback to IP if not authenticated
			}
			return fmt.Sprintf("user:%v", userID)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":       "error",
				"message":      "Too many requests. Please slow down.",
				"retry_after":  int(duration.Seconds()),
			})
		},
		Storage: storage,
	})
}

// SlidingWindowRateLimiter implements a sliding window rate limiter
func SlidingWindowRateLimiter(redisHost string, redisPort int, max int, window time.Duration) fiber.Handler {
	storage := redis.New(redis.Config{
		Host:     redisHost,
		Port:     redisPort,
		Database: 0,
	})

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Create a key with timestamp for sliding window
			now := time.Now().Unix()
			windowStart := now - int64(window.Seconds())
			return fmt.Sprintf("%s:%d", c.IP(), windowStart)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":      "error",
				"message":     "Rate limit exceeded",
				"window":      window.String(),
				"max_requests": max,
			})
		},
		Storage: storage,
	})
}

// EndpointSpecificRateLimiter allows different limits for different endpoints
func EndpointSpecificRateLimiter(redisHost string, redisPort int, limits map[string]limiter.Config) fiber.Handler {
	storage := redis.New(redis.Config{
		Host:     redisHost,
		Port:     redisPort,
		Database: 0,
	})

	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Check if we have a specific limit for this endpoint
		config, exists := limits[path]
		if !exists {
			// Use default limit
			config = limiter.Config{
				Max:        100,
				Expiration: 1 * time.Minute,
			}
		}

		// Set storage and key generator if not provided
		if config.Storage == nil {
			config.Storage = storage
		}
		if config.KeyGenerator == nil {
			config.KeyGenerator = func(c *fiber.Ctx) string {
				return fmt.Sprintf("%s:%s", c.IP(), c.Path())
			}
		}
		if config.LimitReached == nil {
			config.LimitReached = func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"status":  "error",
					"message": "Rate limit exceeded for this endpoint",
				})
			}
		}

		// Apply rate limiting
		return limiter.New(config)(c)
	}
}
