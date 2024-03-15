package api

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	// Response Status Code
	Status int
	// Error Message
	Message any
}

func (e *ApiError) Error() string {
	return fmt.Sprint(e.Message)
}

func buildApiError(status int) func(e error) *ApiError {
	return func(e error) *ApiError {
		return &ApiError{
			Status:  status,
			Message: e,
		}
	}
}

// These errors can be used for known error on every abstraction layer.
var (
	ErrBadRequest     = buildApiError(http.StatusBadRequest)
	ErrUnauthorized   = buildApiError(http.StatusUnauthorized)
	ErrInternalServer = buildApiError(http.StatusInternalServerError)
)
