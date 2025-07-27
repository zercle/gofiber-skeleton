package constant

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInternalServerError   = errors.New("internal server error")
)
