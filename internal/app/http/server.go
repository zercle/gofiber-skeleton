package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
	authRoutes "github.com/zercle/gofiber-skeleton/pkg/domains/auth/api/routes"
	postRoutes "github.com/zercle/gofiber-skeleton/pkg/domains/posts/api/routes"
)

func NewFiberApp() *fiber.App {
	return fiber.New()
}

func RegisterRoutes(app *fiber.App, cfg *config.Config, authRts authRoutes.AuthRoutes, postRts postRoutes.PostRoutes, authMiddleware middleware.AuthMiddleware) {
	authRts.RegisterRoutes(app.Group("/api/v1"), authMiddleware)
	postRts.RegisterRoutes(app.Group("/api/v1"), authMiddleware)
}

func StartServer(app *fiber.App, cfg *config.Config) {
	app.Listen(fmt.Sprintf(":%s", cfg.App.Port))
}
