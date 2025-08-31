package handler

import (
	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes defines and registers all the routes for the user module.
// It maps HTTP endpoints to their corresponding handler functions.
func SetupUserRoutes(router fiber.Router, userHandler *UserHandler) {
	users := router.Group("/users")
	users.Post("/register", userHandler.Register)
}