package errors

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Error parser, parse every error an turn it into ApiError,
// So it can be used to determine what status code should be put on the res headers.
// You can always add your `known error` or make a custom parser for 3rd library/package error.
func Parse(err error) *ApiError {
	// By Error string
	if strings.Contains(err.Error(), "EOF") { // Empty JSON Req body
		return &ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   "EOF, empty json request body",
		}
	}

	// By Error type
	switch e := err.(type) {
	// If the error is instance of ApiErr then no need to do aditional parsing.
	case *ApiError:
		return e
	case validator.ValidationErrors:
		return parseValidationError(e)
	default:
		// It's recommended to log this error, so it can give better insight about the unknown error.
		return &ApiError{
			Status:  500,
			Message: "An unexpected error occurred, try again latter",
		}
	}
}
