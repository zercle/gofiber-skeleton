package userhandler

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers user-related routes
func RegisterRoutes(app *fiber.App) {
	// User authentication & authorization
	userAPI := app.Group("/api/v1")
	userAPI.Post("/register", func(c *fiber.Ctx) error { return c.SendString("User Registration") })
	userAPI.Post("/login", func(c *fiber.Ctx) error { return c.SendString("User Login") })
}