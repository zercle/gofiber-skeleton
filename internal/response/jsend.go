package response

import "github.com/gofiber/fiber/v2"

// JSend response types
const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

// JSendResponse represents a JSend compliant response
type JSendResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// Success sends a JSend success response
func Success(c *fiber.Ctx, httpStatus int, data interface{}) error {
	return c.Status(httpStatus).JSON(JSendResponse{
		Status: StatusSuccess,
		Data:   data,
	})
}

// Fail sends a JSend fail response (client error)
func Fail(c *fiber.Ctx, httpStatus int, data interface{}) error {
	return c.Status(httpStatus).JSON(JSendResponse{
		Status: StatusFail,
		Data:   data,
	})
}

// Error sends a JSend error response (server error)
func Error(c *fiber.Ctx, httpStatus int, message string, code int) error {
	return c.Status(httpStatus).JSON(JSendResponse{
		Status:  StatusError,
		Message: message,
		Code:    code,
	})
}
