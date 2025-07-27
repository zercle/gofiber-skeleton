package http

import (
	"gofiber-skeleton/internal/url/usecase"

	"github.com/gofiber/fiber/v2"
)

// NewURLHandler creates a new URLHandler.
func NewHTTPURLHandler(urlUseCase usecase.URLUseCase) *URLHandler {
	return &URLHandler{urlUseCase: urlUseCase}
}

// URLHandler handles HTTP requests for URLs.
type URLHandler struct {
	urlUseCase usecase.URLUseCase
}

// CreateURL handles URL creation.
func (h *URLHandler) CreateURL(c *fiber.Ctx) error {
	type request struct {
		OriginalURL string `json:"original_url"`
		UserID      string `json:"user_id"` // Optional
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if req.OriginalURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Original URL cannot be empty"})
	}

	url, err := h.urlUseCase.CreateURL(c.Context(), req.OriginalURL, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to create URL"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": url})
}

// GetOriginalURL handles redirection from short code to original URL.
func (h *URLHandler) GetOriginalURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Short code cannot be empty"})
	}

	originalURL, err := h.urlUseCase.GetOriginalURL(c.Context(), shortCode)
	if err != nil {
		if err.Error() == "URL expired" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "URL expired or not found"})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "URL not found"})
	}

	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}
