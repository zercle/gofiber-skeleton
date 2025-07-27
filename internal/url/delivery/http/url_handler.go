package http

import (
	"gofiber-skeleton/internal/url/usecase"

	"github.com/gofiber/fiber/v2"
)

// NewHTTPURLHandler creates a new instance of URLHandler with the provided URL use case.
//
// Parameters:
//   - urlUseCase: The use case interface responsible for URL business logic.
//
// Returns:
//   - *URLHandler: A pointer to the initialized URLHandler instance.
//
// Note:
//   This constructor enables dependency injection of the URL use case.
func NewHTTPURLHandler(urlUseCase usecase.URLUseCase) *URLHandler {
	return &URLHandler{urlUseCase: urlUseCase}
}

// URLHandler handles HTTP requests related to URL operations.
type URLHandler struct {
	urlUseCase usecase.URLUseCase
}

// CreateURLRequest represents the expected JSON payload for creating a new shortened URL.
type CreateURLRequest struct {
	OriginalURL string `json:"original_url"` // The original URL to be shortened.
	UserID      string `json:"user_id"`      // Optional user identifier associated with the URL.
}

// CreateURL handles the creation of a shortened URL.
//
// @Summary Create a new shortened URL
// @Description Accepts an original URL and optional user ID, returning the shortened URL.
// @Tags URLs
// @Accept json
// @Produce json
// @Param request body CreateURLRequest true "URL creation payload"
// @Success 201 {object} map[string]interface{} "Created URL object"
// @Failure 400 {object} map[string]string "Invalid request body or missing original URL"
// @Failure 500 {object} map[string]string "Failed to create URL"
// @Router /api/urls [post]
func (h *URLHandler) CreateURL(c *fiber.Ctx) error {
	var req CreateURLRequest
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

// GetOriginalURL handles redirecting from a short code to the original URL.
//
// @Summary Redirect to original URL
// @Description Redirects the client to the original URL based on the provided short code.
// @Tags URLs
// @Accept json
// @Produce json
// @Param shortCode path string true "Short code for the URL"
// @Success 301 "Redirect to original URL"
// @Failure 400 {object} map[string]string "Short code is empty or invalid"
// @Failure 404 {object} map[string]string "URL not found or expired"
// @Router /{shortCode} [get]
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
