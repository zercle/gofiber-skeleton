package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/handlers"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
)

func SetupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler, cfg *config.Config) {
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	admin := api.Group("/admin")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	authMiddleware := middleware.NewAuthMiddleware(cfg)
	
	auth.Use(authMiddleware)
	auth.Get("/profile", authHandler.GetProfile)
	auth.Put("/profile", authHandler.UpdateProfile)
	auth.Post("/change-password", authHandler.ChangePassword)

	admin.Use(authMiddleware)
	admin.Get("/users", authHandler.ListUsers)
	admin.Get("/users/:id", authHandler.GetUser)
	admin.Post("/users/:id/activate", authHandler.ActivateUser)
	admin.Post("/users/:id/deactivate", authHandler.DeactivateUser)
	admin.Delete("/users/:id", authHandler.DeleteUser)
}