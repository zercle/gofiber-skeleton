package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/handlers"
)

func RegisterAuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler, authMiddleware fiber.Handler) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := auth.Group("", authMiddleware)
	protected.Get("/profile", authHandler.GetProfile)
	protected.Put("/change-password", authHandler.ChangePassword)
}
