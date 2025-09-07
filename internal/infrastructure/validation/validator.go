package validation

import "github.com/go-playground/validator/v10"

// NewValidator creates a new instance of the validator.
func NewValidator() *validator.Validate {
	return validator.New()
}
