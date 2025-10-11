package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground validator
type Validator struct {
	validate *validator.Validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom validators
	_ = v.RegisterValidation("password", validatePassword)
	_ = v.RegisterValidation("username", validateUsername)

	return &Validator{validate: v}
}

// Validate validates a struct
func (v *Validator) Validate(data interface{}) []ValidationError {
	var errors []ValidationError

	err := v.validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: getErrorMessage(err),
			})
		}
	}

	return errors
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
	case "password":
		return fmt.Sprintf("%s must be at least 8 characters and contain uppercase, lowercase, and number", field)
	case "username":
		return fmt.Sprintf("%s must be 3-30 characters and contain only letters, numbers, and underscores", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 30 {
		return false
	}

	return regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
}

// FormatValidationErrors formats validation errors into a string
func FormatValidationErrors(errors []ValidationError) string {
	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}
