package domain

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

const (
	ErrCodeNotFound    = 404
	ErrCodeBadRequest  = 400
	ErrCodeInternal    = 500
	ErrCodeConflict    = 409
	ErrCodeUnauthorized = 401
)

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}