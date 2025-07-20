package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

// GetQRCode godoc
// @Summary Get QR code for a short URL
// @Description Generates and returns a QR code image for a given short URL.
// @Tags urls
// @Produce image/png
// @Param shortCode path string true "Short code of the URL"
// @Success 200 {file} image/png "QR code image"
// @Failure 404 {object} map[string]interface{} "URL not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /{shortCode}/qr [get]
func (h *Handler) GetQRCode(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	longURL, err := h.urlUsecase.GetLongURL(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "data": "URL not found"})
	}

	png, err := qrcode.Encode(longURL, qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate QR code"})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(png)
}
