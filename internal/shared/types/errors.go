package types

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrBadRequest    = errors.New("bad request")
	ErrConflict      = errors.New("conflict")
	ErrInternal      = errors.New("internal server error")
	ErrValidation    = errors.New("validation error")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("expired token")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists   = errors.New("email already exists")
	ErrUserExists    = errors.New("user already exists")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}
	return v[0].Message
}