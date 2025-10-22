package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	apiErr "github.com/prawirdani/golang-restapi/pkg/errors"
)

func parseBodyErr(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	var msg string

	switch {
	case errors.As(err, &syntaxError):
		msg = fmt.Sprintf(
			"Request body contains badly-formed JSON (at position %d)",
			syntaxError.Offset,
		)

	case errors.Is(err, io.ErrUnexpectedEOF):
		msg = "Request body contains badly-formed JSON"

	case errors.As(err, &unmarshalTypeError):
		msg = fmt.Sprintf(
			"Request body contains an invalid value for the %q field (at position %d)",
			unmarshalTypeError.Field,
			unmarshalTypeError.Offset,
		)

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)

	case errors.Is(err, io.EOF):
		msg = "Request body must not be empty"

	default:
		msg = err.Error()
	}
	return apiErr.BadRequest(msg)
}
