package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standardized API response
type Response struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorInfo `json:"error,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
	RequestID string     `json:"request_id,omitempty"`
}

// ErrorInfo represents error information in responses
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       any             `json:"data,omitempty"`
	Error      *ErrorInfo      `json:"error,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
	RequestID  string          `json:"request_id,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Success creates a successful response
func Success(c *fiber.Ctx, statusCode int, message string, data any) error {
	response := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: c.Locals("request_id").(string),
	}

	return c.Status(statusCode).JSON(response)
}

// Error creates an error response
func Error(c *fiber.Ctx, statusCode int, code, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
		RequestID: c.Locals("request_id").(string),
	}

	return c.Status(statusCode).JSON(response)
}

// ErrorWithData creates an error response with additional data
func ErrorWithData(c *fiber.Ctx, statusCode int, code, message string, data any) error {
	response := Response{
		Success: false,
		Message: message,
		Data:    data,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
		RequestID: c.Locals("request_id").(string),
	}

	return c.Status(statusCode).JSON(response)
}

// Paginated creates a paginated response
func Paginated(c *fiber.Ctx, statusCode int, message string, data any, pagination *PaginationInfo) error {
	response := PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Timestamp:  time.Now(),
		RequestID:  c.Locals("request_id").(string),
		Pagination: pagination,
	}

	return c.Status(statusCode).JSON(response)
}

// BadRequest creates a 400 Bad Request response
func BadRequest(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusBadRequest, code, message)
}

// Unauthorized creates a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusUnauthorized, code, message)
}

// Forbidden creates a 403 Forbidden response
func Forbidden(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusForbidden, code, message)
}

// NotFound creates a 404 Not Found response
func NotFound(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusNotFound, code, message)
}

// Conflict creates a 409 Conflict response
func Conflict(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusConflict, code, message)
}

// InternalServerError creates a 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusInternalServerError, code, message)
}

// TooManyRequests creates a 429 Too Many Requests response
func TooManyRequests(c *fiber.Ctx, code, message string) error {
	return Error(c, fiber.StatusTooManyRequests, code, message)
}
