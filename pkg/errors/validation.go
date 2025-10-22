package errors

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *HttpError {
	errors := make(map[string][]string)
	for _, e := range err {
		field := e.Field()

		var msg string
		switch e.Tag() {
		case "required":
			msg = "Field is required"
		case "email":
			msg = "Invalid email format"
		case "min":
			msg = fmt.Sprintf("Must be at least %s characters", e.Param())
		default:
			msg = fmt.Sprintf("Failed on '%s' validation", e.Tag())
		}

		errors[field] = append(errors[field], msg)
	}
	return &HttpError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Invalid request, the provided data does not meet the required format or rules",
		Cause:   errors,
	}
}
