package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// Security creates a security headers middleware
func Security() fiber.Handler {
	return helmet.New()
}
