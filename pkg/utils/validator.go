package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s interface{}) types.ValidationErrors {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors types.ValidationErrors
	
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrs {
			validationErrors = append(validationErrors, types.ValidationError{
				Field:   validationErr.Field(),
				Message: getErrorMessage(validationErr),
			})
		}
	}

	return validationErrors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "len":
		return "Invalid length"
	case "gte":
		return "Value must be greater than or equal to " + err.Param()
	case "lte":
		return "Value must be less than or equal to " + err.Param()
	case "gt":
		return "Value must be greater than " + err.Param()
	case "lt":
		return "Value must be less than " + err.Param()
	case "alpha":
		return "Only alphabetic characters are allowed"
	case "alphanum":
		return "Only alphanumeric characters are allowed"
	case "numeric":
		return "Only numeric characters are allowed"
	case "url":
		return "Invalid URL format"
	case "uri":
		return "Invalid URI format"
	case "uuid":
		return "Invalid UUID format"
	case "uuid4":
		return "Invalid UUID v4 format"
	default:
		return "Invalid value"
	}
}