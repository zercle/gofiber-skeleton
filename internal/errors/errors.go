package errors

import (
	"fmt"
	"net/http"
)

// DomainError represents an application-specific error with context
type DomainError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Cause      error                  `json:"-"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// WithCause adds a cause to the error
func (e *DomainError) WithCause(cause error) *DomainError {
	newErr := *e
	newErr.Cause = cause
	return &newErr
}

// WithContext adds context to the error
func (e *DomainError) WithContext(key string, value interface{}) *DomainError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// New creates a new DomainError
func New(code, message string, httpStatus int) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Context:    make(map[string]interface{}),
	}
}

// Common application errors
var (
	// Authentication errors
	ErrInvalidCredentials = &DomainError{
		Code:       "AUTH_INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenExpired = &DomainError{
		Code:       "AUTH_TOKEN_EXPIRED",
		Message:    "Authentication token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenInvalid = &DomainError{
		Code:       "AUTH_TOKEN_INVALID",
		Message:    "Invalid authentication token",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrUnauthorized = &DomainError{
		Code:       "AUTH_UNAUTHORIZED",
		Message:    "Unauthorized access",
		HTTPStatus: http.StatusUnauthorized,
	}

	// User errors
	ErrUserNotFound = &DomainError{
		Code:       "USER_NOT_FOUND",
		Message:    "User not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUserAlreadyExists = &DomainError{
		Code:       "USER_ALREADY_EXISTS",
		Message:    "User with this email already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrUserInactive = &DomainError{
		Code:       "USER_INACTIVE",
		Message:    "User account is inactive",
		HTTPStatus: http.StatusForbidden,
	}

	// Validation errors
	ErrValidationFailed = &DomainError{
		Code:       "VALIDATION_FAILED",
		Message:    "Validation failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidInput = &DomainError{
		Code:       "INVALID_INPUT",
		Message:    "Invalid input provided",
		HTTPStatus: http.StatusBadRequest,
	}

	// Database errors
	ErrDatabaseOperation = &DomainError{
		Code:       "DATABASE_ERROR",
		Message:    "Database operation failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrRecordNotFound = &DomainError{
		Code:       "RECORD_NOT_FOUND",
		Message:    "Record not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrDuplicateKey = &DomainError{
		Code:       "DUPLICATE_KEY",
		Message:    "Duplicate key constraint violation",
		HTTPStatus: http.StatusConflict,
	}

	// General errors
	ErrInternal = &DomainError{
		Code:       "INTERNAL_ERROR",
		Message:    "An internal error occurred",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrNotFound = &DomainError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrForbidden = &DomainError{
		Code:       "FORBIDDEN",
		Message:    "Access forbidden",
		HTTPStatus: http.StatusForbidden,
	}

	ErrRateLimitExceeded = &DomainError{
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    "Rate limit exceeded",
		HTTPStatus: http.StatusTooManyRequests,
	}

	// External service errors
	ErrExternalService = &DomainError{
		Code:       "EXTERNAL_SERVICE_ERROR",
		Message:    "External service error",
		HTTPStatus: http.StatusBadGateway,
	}

	ErrServiceUnavailable = &DomainError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Service temporarily unavailable",
		HTTPStatus: http.StatusServiceUnavailable,
	}
)

// IsDomainError checks if an error is a DomainError
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// GetErrorCode returns the error code for an error
func GetErrorCode(err error) string {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Code
	}
	return "UNKNOWN_ERROR"
}
