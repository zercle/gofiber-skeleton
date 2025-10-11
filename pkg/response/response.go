package response

import "github.com/gofiber/fiber/v2"

// Response represents a standard API response using JSend format
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents error details
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Status     string         `json:"status"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// Success sends a success response
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Resource created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

// NoContent sends a no content response
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// BadRequest sends a bad request error response
func BadRequest(c *fiber.Ctx, message string, details ...string) error {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: detail,
		},
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *fiber.Ctx, message ...string) error {
	msg := "Unauthorized access"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "UNAUTHORIZED",
			Message: msg,
		},
	})
}

// Forbidden sends a forbidden error response
func Forbidden(c *fiber.Ctx, message ...string) error {
	msg := "Access forbidden"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusForbidden).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "FORBIDDEN",
			Message: msg,
		},
	})
}

// NotFound sends a not found error response
func NotFound(c *fiber.Ctx, resource ...string) error {
	msg := "Resource not found"
	if len(resource) > 0 {
		msg = resource[0] + " not found"
	}

	return c.Status(fiber.StatusNotFound).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "NOT_FOUND",
			Message: msg,
		},
	})
}

// Conflict sends a conflict error response
func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

// UnprocessableEntity sends an unprocessable entity error response
func UnprocessableEntity(c *fiber.Ctx, message string, details ...string) error {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "UNPROCESSABLE_ENTITY",
			Message: message,
			Details: detail,
		},
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(c *fiber.Ctx, message ...string) error {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: msg,
		},
	})
}

// Paginated sends a paginated response
func Paginated(c *fiber.Ctx, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Status:     "success",
		Data:       data,
		Pagination: meta,
	})
}

// Custom sends a custom response with specific status code
func Custom(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorWithCode sends an error response with custom error code
func ErrorWithCode(c *fiber.Ctx, statusCode int, errorCode string, message string, details ...string) error {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	return c.Status(statusCode).JSON(Response{
		Status: "error",
		Error: &Error{
			Code:    errorCode,
			Message: message,
			Details: detail,
		},
	})
}
