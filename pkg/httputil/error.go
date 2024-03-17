package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Use these error wrappers for known errors to have precise response status codes, can be used on any abstraction layer.
// You can either pass string, error object or anything as the error cause message.
var (
	ErrBadRequest       = buildApiError(http.StatusBadRequest)
	ErrConflict         = buildApiError(http.StatusConflict)
	ErrNotFound         = buildApiError(http.StatusNotFound)
	ErrUnauthorized     = buildApiError(http.StatusUnauthorized)
	ErrInternalServer   = buildApiError(http.StatusInternalServerError)
	ErrMethodNotAllowed = buildApiError(http.StatusMethodNotAllowed)
)

type ApiError struct {
	// Response Status Code
	Status int
	// Error Cause
	Cause interface{}
}

// Return ApiErr in string format
func (e *ApiError) Error() string {
	return fmt.Sprintf("status: %v cause: %v", e.Status, e.Cause)
}

func buildApiError(status int) func(err interface{}) *ApiError {
	return func(e interface{}) *ApiError {
		return &ApiError{
			Status: status,
			Cause:  e,
		}
	}
}

// Error parser, parse every error an turn it into ApiError,
// So it can be used to determine what status code should be put on the res headers.
// You can always add your `known error` or make a custom parser for 3rd library/package error.
func parseError(err error) *ApiError {
	// By Error string
	switch {
	case strings.Contains(err.Error(), "EOF"): // Empty JSON Req body
		return ErrBadRequest("Empty json request body")
	}

	// By Error type
	switch e := err.(type) {
	// If the error is instance of ApiErr then no need to do aditional parsing.
	case *ApiError:
		return e
	case validator.ValidationErrors:
		return parseValidationError(e)
	case *json.UnmarshalTypeError:
		return parseJsonUnmarshalTypeError(e)
	default:
		return ErrInternalServer(err.Error())
	}

}

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *ApiError {
	// Validation error mapped into a map, so the response will look like "field":"the error"
	errors := make(map[string]interface{})
	for _, errField := range err {
		field := strings.ToLower(errField.Field())
		switch errField.Tag() {
		case "required":
			errors[field] = fmt.Sprintf("%s field is required", field)
		case "email":
			errors[field] = "Invalid email format"
		default:
			errors[field] = errField.Error()
		}
	}
	return ErrBadRequest(errors)
}

// JSON Unmarshal mismatch type error
func parseJsonUnmarshalTypeError(err *json.UnmarshalTypeError) *ApiError {
	if strings.Contains(err.Error(), "unmarshal") {
		return ErrBadRequest(fmt.Sprintf("Type mismatch at %s, Expected type %s, Got %s", err.Field, err.Type, err.Value))
	}
	return ErrBadRequest(err.Error())
}
