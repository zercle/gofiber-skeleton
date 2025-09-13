package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/api/handlers"
)

type AuthRoutes struct {
	authHandler handlers.AuthHandler
}

func NewAuthRoutes(authHandler handlers.AuthHandler) AuthRoutes {
	return AuthRoutes{authHandler: authHandler}
}

func (r *AuthRoutes) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/login", r.authHandler.Login)
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/refresh", r.authHandler.RefreshToken)

	// Protected routes
	protected := auth.Group("", authMiddleware)
	protected.Get("/profile", r.authHandler.GetProfile)
	protected.Put("/change-password", r.authHandler.ChangePassword)
}