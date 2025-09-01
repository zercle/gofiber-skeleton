package handler

import (
	"github.com/gofiber/fiber/v2"
)

// InitUserRoutes defines and registers all the routes for the user module.
// It maps HTTP endpoints to their corresponding handler functions.
func (h *UserHandler) InitUserRoutes(app *fiber.App, jwtMiddleware fiber.Handler) {
	users := app.Group("/api/v1/users")
	users.Post("/", h.Register)

	// Protected routes
	users.Use(jwtMiddleware)
	users.Get("/", h.GetAll)
	users.Get("/:id", h.GetByID)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}
