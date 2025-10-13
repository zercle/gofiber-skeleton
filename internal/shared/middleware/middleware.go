package middleware

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

func NewRecovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			fmt.Printf("panic recovered: %v\n%s", e, string(c.Request().Header.RawHeaders()))
		},
	})
}

func NewLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:       "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat:   "2006-01-02 15:04:05",
		TimeZone:     "UTC",
		Output:       os.Stdout,
		DisableColors: false,
	})
}

func NewCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		ExposeHeaders:    "X-Total-Count,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           86400,
	})
}

func NewRequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
	})
}

func NewTimeout(timeoutDuration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeoutDuration)
		defer cancel()

		c.SetUserContext(ctx)
		return c.Next()
	}
}

func NewRateLimiting() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-RateLimit-Limit", "100")
		c.Set("X-RateLimit-Remaining", "99")
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))

		return c.Next()
	}
}

func NewSecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Content-Security-Policy", "default-src 'self'")

		return c.Next()
	}
}

func NewHealthCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "gofiber-skeleton",
			"version":   "1.0.0",
		})
	}
}