package response

import "github.com/gofiber/fiber/v2"

// JSendSuccess defines the structure for a successful JSend response.
type JSendSuccess struct {
	Status string      `json:"status"` // Always "success"
	Data   interface{} `json:"data"`   // The data payload
}

// JSendFail defines the structure for a failed JSend response.
type JSendFail struct {
	Status string      `json:"status"` // Always "fail"
	Data   interface{} `json:"data"`   // Details about the failure
}

// JSendError defines the structure for an error JSend response.
type JSendError struct {
	Status  string      `json:"status"`         // Always "error"
	Message string      `json:"message"`        // A meaningful, human-readable error message
	Code    interface{} `json:"code,omitempty"` // A unique error code, optional
	Data    interface{} `json:"data,omitempty"` // Additional error data, optional
}

// Success sends a JSend success response.
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(JSendSuccess{
		Status: "success",
		Data:   data,
	})
}

// Fail sends a JSend fail response.
func Fail(c *fiber.Ctx, message string, data interface{}, statusCode int) error {
	c.Status(statusCode)
	return c.JSON(JSendFail{
		Status: "fail",
		Data: fiber.Map{
			"message": message,
			"details": data,
		},
	})
}

// Error sends a JSend error response.
func Error(c *fiber.Ctx, message string, code interface{}, data interface{}, statusCode int) error {
	c.Status(statusCode)
	return c.JSON(JSendError{
		Status:  "error",
		Message: message,
		Code:    code,
		Data:    data,
	})
}
