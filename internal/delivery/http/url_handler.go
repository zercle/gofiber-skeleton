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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if req.OriginalURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Original URL cannot be empty"})
	}


	user, ok := c.Locals("user").(*jwt.Token)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid user ID in token"})
	}

	url, err := h.urlUseCase.CreateShortURL(c.Context(), req.OriginalURL, userID, req.CustomShort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to create short URL"})
	}
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Short URL created successfully",
		"data": fiber.Map{
			"short_code": url.ShortCode,
		},
	})
}

// Redirect redirects a short URL to the original URL.
func (h *URLHandler) Redirect(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	originalURL, err := h.urlUseCase.GetOriginalURL(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "data": "Short URL not found"})
	}

	// Basic validation to prevent open redirect. A more robust solution would involve a whitelist.
	if !isValidURL(originalURL) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid original URL"})
	}

	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}

// isValidURL performs a basic check to ensure the URL is safe for redirection.
func isValidURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}

// GetQRCode generates a QR code for a short URL.
func (h *URLHandler) GetQRCode(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	qrCode, err := h.urlUseCase.GenerateQRCode(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate QR code"})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(qrCode)
}
