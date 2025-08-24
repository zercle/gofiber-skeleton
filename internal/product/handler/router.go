package producthandler

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers product-related routes
func RegisterRoutes(app *fiber.App) {
	// Product management routes
	productAPI := app.Group("/api/v1/products")
	productAPI.Post("/", func(c *fiber.Ctx) error { return c.SendString("Create Product") })
	productAPI.Put("/:id", func(c *fiber.Ctx) error { return c.SendString("Update Product") })
	productAPI.Delete("/:id", func(c *fiber.Ctx) error { return c.SendString("Delete Product") })
	productAPI.Get("/", func(c *fiber.Ctx) error { return c.SendString("Get All Products") })
	productAPI.Get("/:id", func(c *fiber.Ctx) error { return c.SendString("Get Product by ID") })
}