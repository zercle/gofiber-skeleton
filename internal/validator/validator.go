package validator

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom tag name function to use json tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: v}
}

// Validate validates a struct
func (v *Validator) Validate(data interface{}) []ValidationError {
	var validationErrors []ValidationError

	err := v.validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Message: getErrorMessage(err),
			})
		}
	}

	return validationErrors
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// getErrorMessage returns a user-friendly error message
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "eqfield":
		return fmt.Sprintf("%s must equal %s", field, err.Param())
	case "nefield":
		return fmt.Sprintf("%s must not equal %s", field, err.Param())
	default:
		return fmt.Sprintf("%s failed %s validation", field, tag)
	}
}
