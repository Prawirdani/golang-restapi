package errors

import (
	"net/http"
)

// Use these error wrappers for known errors to have precise response status codes, can be used on any abstraction layer.
// It will only set the message in ErrorResponse, if you want to provide details in the ErrorResponse you should create ApiError object manually.
var (
	BadRequest       = build(http.StatusBadRequest)
	Conflict         = build(http.StatusConflict)
	NotFound         = build(http.StatusNotFound)
	Unauthorized     = build(http.StatusUnauthorized)
	MethodNotAllowed = build(http.StatusMethodNotAllowed)
	Forbidden        = build(http.StatusForbidden)
	InternalServer   = build(http.StatusInternalServerError)
)

type ApiError struct {
	Status  int
	Message string
	Cause   interface{}
}

// Return ApiErr in string format
func (e *ApiError) Error() string {
	return e.Message
}

func build(status int) func(msg string) *ApiError {
	return func(m string) *ApiError {
		return &ApiError{
			Status:  status,
			Message: m,
		}
	}
}
