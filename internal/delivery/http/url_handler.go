package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type createURLRequest struct {
	LongURL   string `json:"long_url"`
	CustomShort string `json:"custom_short,omitempty"`
}

// CreateShortURL godoc
// @Summary Create a short URL
// @Description Creates a new short URL, optionally with a custom short code.
// @Tags urls
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param url body createURLRequest true "URL creation details"
// @Success 201 {object} entities.URL "Short URL created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /urls [post]
func (h *Handler) CreateShortURL(c *fiber.Ctx) error {
	var req createURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "data": err.Error()})
	}

	var userID *string
	user := c.Locals("user")
	if user != nil {
		t := user.(*jwt.Token)
		claims := t.Claims.(jwt.MapClaims)
		uid := claims["user_id"].(string)
		userID = &uid
	}

	url, err := h.urlUsecase.CreateShortURL(c.Context(), userID, req.LongURL, req.CustomShort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": url})
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects from a short code to its original long URL.
// @Tags urls
// @Accept json
// @Produce html
// @Param shortCode path string true "Short code of the URL"
// @Success 301 "Redirect to original URL"
// @Failure 404 {object} map[string]interface{} "URL not found"
// @Router /{shortCode} [get]
func (h *Handler) Redirect(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	longURL, err := h.urlUsecase.GetLongURL(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "data": "URL not found"})
	}

	return c.Redirect(longURL, fiber.StatusMovedPermanently)
}

// GetUserURLs godoc
// @Summary Get all URLs for a user
// @Description Retrieves all short URLs created by the authenticated user.
// @Tags urls
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} entities.URL "List of URLs"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /urls [get]
func (h *Handler) GetUserURLs(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	urls, err := h.urlUsecase.GetUserURLs(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": urls})
}

// DeleteURL godoc
// @Summary Delete a URL
// @Description Deletes a short URL by its ID for the authenticated user.
// @Tags urls
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "URL ID"
// @Success 204 "URL deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /urls/{id} [delete]
func (h *Handler) DeleteURL(c *fiber.Ctx) error {
	id := c.Params("id")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	if err := h.urlUsecase.DeleteURL(c.Context(), id, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
