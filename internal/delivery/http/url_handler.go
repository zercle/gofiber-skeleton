package http

import (
	"gofiber-skeleton/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// NewURLHandler creates a new URLHandler.
func NewURLHandler(urlUseCase usecases.URLUseCase) *URLHandler {
	return &URLHandler{urlUseCase: urlUseCase}
}

// URLHandler handles HTTP requests for URLs.
type URLHandler struct {
	urlUseCase usecases.URLUseCase
}

// CreateShortURL handles the creation of a short URL.
func (h *URLHandler) CreateShortURL(c *fiber.Ctx) error {
	type request struct {
		OriginalURL string `json:"original_url"`
		CustomShort string `json:"custom_short,omitempty"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "data": err.Error()})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, _ := uuid.Parse(claims["sub"].(string))

	url, err := h.urlUseCase.CreateShortURL(c.Context(), req.OriginalURL, userID, req.CustomShort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": url})
}

// Redirect redirects a short URL to the original URL.
func (h *URLHandler) Redirect(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	originalURL, err := h.urlUseCase.GetOriginalURL(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "data": "Short URL not found"})
	}

	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}

// GetQRCode generates a QR code for a short URL.
func (h *URLHandler) GetQRCode(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	qrCode, err := h.urlUseCase.GenerateQRCode(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(qrCode)
}
