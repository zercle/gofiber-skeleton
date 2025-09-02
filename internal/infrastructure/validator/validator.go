package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/pkg/response"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
	Message     string
}

var validate = validator.New()

func ValidateStruct(data any) []ErrorResponse {
	var errors []ErrorResponse

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Message = getErrorMessage(err)
			errors = append(errors, element)
		}
	}
	return errors
}

func getErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, fe.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be a valid number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s is not valid", field)
	}
}

func HandleValidationErrors(c *fiber.Ctx, errors []ErrorResponse) error {
	errorMap := make(map[string]string)
	for _, err := range errors {
		field := strings.ToLower(strings.ReplaceAll(err.FailedField, ".", "_"))
		errorMap[field] = err.Message
	}

	resp := response.Fail(map[string]any{
		"validation_errors": errorMap,
	})

	return c.Status(fiber.StatusUnprocessableEntity).JSON(resp)
}
