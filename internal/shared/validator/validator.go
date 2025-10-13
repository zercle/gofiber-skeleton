package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(s any) error {
	return validate.Struct(s)
}

func GetValidator() *validator.Validate {
	return validate
}
