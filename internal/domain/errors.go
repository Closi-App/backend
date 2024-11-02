package domain

import "fmt"

var (
	ErrInternalServerError = NewError("ERR_INTERNAL_SERVER_ERROR", "internal server error")
	ErrBadRequest          = NewError("ERR_BAD_REQUEST", "bad request")
	ErrUnauthorized        = NewError("ERR_UNAUTHORIZED", "unauthorized access")
)

type Error struct {
	Code    string
	Message string
}

func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
