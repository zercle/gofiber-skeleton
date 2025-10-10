package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Recovery returns recovery middleware to handle panics
func Recovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}
