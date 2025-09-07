package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/handlers"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/repositories"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/usecases"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
)

// RegisterRoutes registers all auth module routes.
func RegisterRoutes(router fiber.Router, db *database.Database, cfg *config.Config, authMiddleware fiber.Handler) {
	// Initialize repository, usecase, and handler
	userRepo := repositories.NewUserRepository(db)
	authUseCase := usecases.NewAuthUseCase(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authUseCase)

	auth := router.Group("/auth")

	// Public routes
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	auth.Get("/profile", authMiddleware, authHandler.GetProfile)
	auth.Put("/change-password", authMiddleware, authHandler.ChangePassword)
}
