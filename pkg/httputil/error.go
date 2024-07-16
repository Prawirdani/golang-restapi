package httputil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Use these error wrappers for known errors to have precise response status codes, can be used on any abstraction layer.
// It will only set the message in ErrorResponse, if you want to provide details in the ErrorResponse you should create custom ApiError object.
var (
	ErrBadRequest       = buildApiError(http.StatusBadRequest)
	ErrConflict         = buildApiError(http.StatusConflict)
	ErrNotFound         = buildApiError(http.StatusNotFound)
	ErrUnauthorized     = buildApiError(http.StatusUnauthorized)
	ErrMethodNotAllowed = buildApiError(http.StatusMethodNotAllowed)
	ErrInternalServer   = buildApiError(http.StatusInternalServerError)
)

type apiError struct {
	status  int
	message string
	cause   interface{}
}

// Return ApiErr in string format
func (e *apiError) Error() string {
	return e.message
}

func buildApiError(status int) func(msg string) *apiError {
	return func(m string) *apiError {
		return &apiError{
			status:  status,
			message: m,
		}
	}
}

// Error parser, parse every error an turn it into ApiError,
// So it can be used to determine what status code should be put on the res headers.
// You can always add your `known error` or make a custom parser for 3rd library/package error.
func parseErrors(err error) *apiError {
	// By Error string
	if strings.Contains(err.Error(), "EOF") { // Empty JSON Req body
		return &apiError{
			status:  http.StatusBadRequest,
			message: "Invalid request body",
			cause:   "EOF, empty json request body",
		}
	}

	// By Error type
	switch e := err.(type) {
	// If the error is instance of ApiErr then no need to do aditional parsing.
	case *apiError:
		return e
	case validator.ValidationErrors:
		return parseValidationError(e)
	case *json.UnmarshalTypeError:
		return parseJsonUnmarshalTypeError(e)
	case *json.SyntaxError:
		return parseJsonSyntaxError(e)
	default:
		// Log the unknown error
		errReflectType := reflect.TypeOf(err) // Determine the reflect type of the error for easier examination
		slog.Error("Unknown ERROR", slog.Any("cause", err), slog.String("reflectType", errReflectType.String()))
		return &apiError{
			status:  500,
			message: "An unexpected error occurred, try again latter",
			cause:   err.Error(),
		}

	}

}

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *apiError {
	// Validation error mapped into a map, so the response will look like "field":"the error"
	errors := make(map[string]interface{})
	for _, errField := range err {
		field := strings.ToLower(errField.Field())
		switch errField.Tag() {
		case "required":
			errors[field] = "Field is required"
		case "email":
			errors[field] = "Invalid email format"
		case "min":
			if field == "password" {
				errors[field] = "Must be at least 6 characters long"
			} else {
				errors[field] = fmt.Sprintf("Must be at least %s characters long", errField.Param())
			}
		default:
			errors[field] = errField.Error()
		}
	}
	return &apiError{
		status:  http.StatusUnprocessableEntity,
		message: "Invalid request, the provided data does not meet the required format or rules",
		cause:   errors,
	}
}

// JSON Unmarshal mismatch type error
func parseJsonUnmarshalTypeError(err *json.UnmarshalTypeError) *apiError {
	e := &apiError{
		status:  http.StatusBadRequest,
		message: "Invalid request body, type error",
		cause:   err.Error(),
	}
	if strings.Contains(err.Error(), "unmarshal") {
		e.cause = fmt.Sprintf("Type mismatch at %s field, expected type %s, got %s", err.Field, err.Type, err.Value)
	}
	return e
}

func parseJsonSyntaxError(err *json.SyntaxError) *apiError {
	return &apiError{
		status:  http.StatusBadRequest,
		message: "Invalid request body, syntax error",
		cause:   err.Error(),
	}

}
