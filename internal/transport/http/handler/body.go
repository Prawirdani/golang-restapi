package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// MaxBodySize maximum read size from request body
const MaxBodySize = 32 << 20 // 32 MB

// Body is json response body
type Body struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// MarshalJSON implements json.Marshaller to prevent using pointer on Body fields while preserving the field on the body with nullable ability
func (b Body) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"data": b.Data,
	}

	if b.Message != "" {
		m["message"] = b.Message
	} else {
		m["message"] = nil
	}

	return json.Marshal(m)
}

// JSONRequestBody is an interface for request body that needs to be validated and sanitized when binding to struct on handler function
// It is recommended to implement this interface on every request body model struct
// Should use this alongside with BindValidate function from binder.go
type JSONRequestBody interface {
	// Validate is a method to validate request body to ensure all required fields are provided and match the constraints
	Validate() error
	// Sanitize is a method to sanitize request body to ensure all fields are in the correct format and clean
	Sanitize() error
}

// eTag generate strong etag from given data
func eTag(data any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	h := sha256.Sum256(b)
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(h[:]))
}
