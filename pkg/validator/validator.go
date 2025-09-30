package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/response"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

// ValidateRequest validates a request struct and sends appropriate error response
func ValidateRequest(c *fiber.Ctx, data interface{}) error {
	if err := ValidateStruct(data); err != nil {
		validationErrors := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				field := strings.ToLower(fe.Field())
				validationErrors[field] = formatValidationError(fe)
			}
		}

		return response.Fail(c, fiber.StatusBadRequest, fiber.Map{
			"validation_errors": validationErrors,
		})
	}
	return nil
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	default:
		return fmt.Sprintf("%s failed %s validation", fe.Field(), fe.Tag())
	}
}