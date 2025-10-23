package response

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/errors"
)

type Body struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"-"`
	Errors  any    `json:"errors,omitempty"`
}

// Send is a function to send json response to client
// It uses option pattern to accepts multiple options to customize the response
func Send(w http.ResponseWriter, r *http.Request, opts ...Option) error {
	res := Body{
		Status: http.StatusOK, // Default
	}
	for _, opt := range opts {
		opt(&res)
	}

	etag := eTag(res)
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	w.Header().Set("ETag", etag)
	return writeJSON(w, res.Status, res)
}

// Response writer for handling error
func HandleError(w http.ResponseWriter, err error) {
	e := errors.Parse(err)
	response := Body{
		Message: e.Message,
		Errors:  e.Cause,
	}

	writeErr := writeJSON(w, e.Status, response)
	if writeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func eTag(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "error"
	}
	// Compute the SHA-256 hash of the JSON representation
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf(`"%x"`, hash)
}
