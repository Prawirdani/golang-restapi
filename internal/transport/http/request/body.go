package request

import (
	"encoding/json"
	"net/http"

	apiErr "github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

// Body is an interface for request body that needs to be validated and sanitized when binding to struct on handler function
// It is recommended to implement this interface on every request body model struct
// Should use this alongside with BindValidate function from binder.go
type Body interface {
	// Validate is a method to validate request body to ensure all required fields are provided and match the constraints
	Validate() error
	// Sanitize is a method to sanitize request body to ensure all fields are in the correct format and clean
	Sanitize() error
}

// BindValidate is a helper function to bind and validate json request body, requires a struct that implements RequestBody interface
func BindValidate(r *http.Request, body any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return apiErr.BadRequest("Content-Type must be application/json")
	}

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(body); err != nil {
		return parseJsonBodyErr(err)
	}

	// If Implement Body interfaces
	if b, ok := body.(Body); ok {
		if err := b.Sanitize(); err != nil {
			return err
		}

		return b.Validate()
	}

	// Do Manual validation
	return validator.Struct(body)
}
