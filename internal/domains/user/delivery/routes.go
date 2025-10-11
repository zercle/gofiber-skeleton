package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

// RegisterRoutes registers user routes
func RegisterRoutes(router fiber.Router, handler *UserHandler, cfg *config.Config) {
	// Public routes
	auth := router.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Protected routes
	users := router.Group("/users", middleware.AuthMiddleware(cfg))
	users.Get("/profile", handler.GetProfile)
	users.Put("/profile", handler.UpdateProfile)
	users.Put("/password", handler.ChangePassword)
	users.Get("/:id", handler.GetUserByID)
	users.Get("", handler.ListUsers)
	users.Post("/:id/deactivate", handler.DeactivateUser)
	users.Post("/:id/activate", handler.ActivateUser)
	users.Delete("/:id", handler.DeleteUser)
}
