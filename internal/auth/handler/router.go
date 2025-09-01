package handler

import "github.com/gofiber/fiber/v2"

func (h *AuthHandler) InitAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", h.Login)
}