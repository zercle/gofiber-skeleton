package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/middleware"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	appMiddleware "github.com/zercle/gofiber-skeleton/internal/middleware"
)

// SetupUserRoutes sets up user-related routes
func SetupUserRoutes(router fiber.Router, authUsecase usecase.AuthUsecase) {
	authHandler := NewAuthHandler(authUsecase)

	// Public routes (no authentication required)
	auth := router.Group("/auth")
	auth.Post("/register", appMiddleware.AuthRateLimit(), authHandler.Register)
	auth.Post("/login", appMiddleware.AuthRateLimit(), authHandler.Login)

	// Protected routes (authentication required)
	users := router.Group("/users")
	users.Use(middleware.JWTAuth(authUsecase))
	users.Get("/me", authHandler.GetProfile)
	users.Put("/me", authHandler.UpdateProfile)
	users.Put("/me/password", authHandler.ChangePassword)
}
