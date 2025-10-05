package errors

import (
	"fmt"
)

// AppError represents an application-specific error
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	CodeNotFound          = "NOT_FOUND"
	CodeUnauthorized      = "UNAUTHORIZED"
	CodeForbidden         = "FORBIDDEN"
	CodeBadRequest        = "BAD_REQUEST"
	CodeInternalError     = "INTERNAL_ERROR"
	CodeValidationError   = "VALIDATION_ERROR"
	CodeDuplicateError    = "DUPLICATE_ERROR"
	CodeDatabaseError     = "DATABASE_ERROR"
	CodeConfigError       = "CONFIG_ERROR"
)

// NewAppError creates a new AppError
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NotFound creates a not found error
func NotFound(message string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: message,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: message,
	}
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Code:    CodeForbidden,
		Message: message,
	}
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return &AppError{
		Code:    CodeBadRequest,
		Message: message,
	}
}

// InternalError creates an internal error
func InternalError(message string, err error) *AppError {
	return &AppError{
		Code:    CodeInternalError,
		Message: message,
		Err:     err,
	}
}

// ValidationError creates a validation error
func ValidationError(message string) *AppError {
	return &AppError{
		Code:    CodeValidationError,
		Message: message,
	}
}

// DuplicateError creates a duplicate error
func DuplicateError(message string) *AppError {
	return &AppError{
		Code:    CodeDuplicateError,
		Message: message,
	}
}

// DatabaseError creates a database error
func DatabaseError(message string, err error) *AppError {
	return &AppError{
		Code:    CodeDatabaseError,
		Message: message,
		Err:     err,
	}
}

// Is checks if an error is of a specific type
func Is(err error, code string) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}
