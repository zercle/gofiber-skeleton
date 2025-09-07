package jsend

import "github.com/gofiber/fiber/v2"

// JSendResponse represents the JSend specification for API responses.
type JSendResponse struct {
	Status  string      `json:"status"`            // "success", "fail", or "error"
	Data    interface{} `json:"data,omitempty"`    // Data for success/fail
	Message string      `json:"message,omitempty"` // Error message for error
	Code    int         `json:"code,omitempty"`    // Optional error code
}

// SendSuccess sends a JSend success response.
func SendSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(JSendResponse{
		Status: "success",
		Data:   data,
	})
}

// SendFail sends a JSend fail response (client-side error).
func SendFail(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(JSendResponse{
		Status: "fail",
		Data:   data,
	})
}

// SendError sends a JSend error response (server-side error).
func SendError(c *fiber.Ctx, statusCode int, message string, code ...int) error {
	resp := JSendResponse{
		Status:  "error",
		Message: message,
	}
	if len(code) > 0 {
		resp.Code = code[0]
	}
	return c.Status(statusCode).JSON(resp)
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// SendValidationError sends validation error response
func SendValidationError(c *fiber.Ctx, errors []ValidationError) error {
	return SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
		"message": "Validation failed",
		"errors":  errors,
	})
}

// SendUnauthorized sends a 401 Unauthorized JSend fail response.
func SendUnauthorized(c *fiber.Ctx, data interface{}) error {
	return SendFail(c, fiber.StatusUnauthorized, data)
}

// SendBadRequest sends a 400 Bad Request JSend fail response.
func SendBadRequest(c *fiber.Ctx, data interface{}) error {
	return SendFail(c, fiber.StatusBadRequest, data)
}

// SendNotFound sends a 404 Not Found JSend fail response.
func SendNotFound(c *fiber.Ctx, data interface{}) error {
	return SendFail(c, fiber.StatusNotFound, data)
}

// SendForbidden sends a 403 Forbidden JSend fail response.
func SendForbidden(c *fiber.Ctx, data interface{}) error {
	return SendFail(c, fiber.StatusForbidden, data)
}

// SendInternalServerError sends a 500 Internal Server Error JSend error response.
func SendInternalServerError(c *fiber.Ctx, message string, code ...int) error {
	return SendError(c, fiber.StatusInternalServerError, message, code...)
}
