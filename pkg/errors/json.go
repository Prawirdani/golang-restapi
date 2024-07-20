package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// JSON Unmarshal mismatch type error
func parseJsonUnmarshalTypeError(err *json.UnmarshalTypeError) *ApiError {
	e := &ApiError{
		Status:  http.StatusBadRequest,
		Message: "Invalid request body, type error",
		Cause:   err.Error(),
	}
	if strings.Contains(err.Error(), "unmarshal") {
		e.Cause = fmt.Sprintf("Type mismatch at %s field, expected type %s, got %s", err.Field, err.Type, err.Value)
	}
	return e
}

func parseJsonSyntaxError(err *json.SyntaxError) *ApiError {
	return &ApiError{
		Status:  http.StatusBadRequest,
		Message: "Invalid request body, syntax error",
		Cause:   err.Error(),
	}

}
