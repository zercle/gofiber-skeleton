package http

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-skeleton/internal/usecases"
)

type Handler struct {
	userUsecase *usecases.UserUsecase
	urlUsecase  *usecases.URLUsecase
}

func NewHandler(userUsecase *usecases.UserUsecase, urlUsecase *usecases.URLUsecase) *Handler {
	return &Handler{userUsecase: userUsecase, urlUsecase: urlUsecase}
}

func (h *Handler) Register(app *fiber.App) {
	api := app.Group("/api")

	// User routes
	api.Post("/register", h.RegisterUser)
	api.Post("/login", h.Login)

	// URL routes
	api.Post("/urls", AuthMiddleware, h.CreateShortURL)
	app.Get("/:shortCode", h.Redirect)

	// Authenticated routes
	auth := api.Group("/", AuthMiddleware)
	auth.Get("/urls", h.GetUserURLs)
	auth.Delete("/urls/:id", h.DeleteURL)

	// QR code route
	app.Get("/:shortCode/qr", h.GetQRCode)

	// Health check route
	app.Get("/health", h.HealthCheck)
}

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
