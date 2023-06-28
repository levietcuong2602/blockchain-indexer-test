package httperr

import (
	"net/http"
)

var (
	ErrBadRequest        = NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	ErrInternalServer    = NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	ErrCollectionExisted = NewError(http.StatusInternalServerError, "Collection existed")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetStatusCode() int {
	return e.Code
}
