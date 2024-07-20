package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *ApiError {
	// Validation error mapped into a map, so the response will look like "field":"the error"
	errors := make(map[string]interface{})
	for _, errField := range err {
		field := strings.ToLower(errField.Field()[0:1]) + errField.Field()[1:]
		switch errField.Tag() {
		case "required":
			errors[field] = "Field is required"
		case "email":
			errors[field] = "Invalid email format"
		case "min":
			errors[field] = fmt.Sprintf("Must be at least %s characters long", errField.Param())
		case "eqfield":
			errors[field] = fmt.Sprintf("Must be the same as %s", strings.ToLower(errField.Param()[0:1])+errField.Param()[1:])
		default:
			errors[field] = errField.Error()
		}
	}
	return &ApiError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Invalid request, the provided data does not meet the required format or rules",
		Cause:   errors,
	}
}
