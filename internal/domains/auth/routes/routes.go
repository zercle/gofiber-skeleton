package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/handlers"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/usecases"
)

// RegisterRoutes registers all auth module routes.
func RegisterRoutes(router fiber.Router, authUseCase usecases.AuthUseCase, authHandler handlers.AuthHandler, authMiddleware fiber.Handler) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	auth.Get("/profile", authMiddleware, authHandler.GetProfile)
	auth.Put("/change-password", authMiddleware, authHandler.ChangePassword)
}
