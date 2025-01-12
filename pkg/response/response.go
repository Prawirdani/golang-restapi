package response

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseBody struct {
	Data    any         `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"-"`
	Errors  interface{} `json:"errors"`
}

// Send is a function to send json response to client
// It uses option pattern to accepts multiple options to customize the response
func Send(w http.ResponseWriter, r *http.Request, opts ...Option) error {
	res := ResponseBody{
		Status: http.StatusOK, // Default
	}
	for _, opt := range opts {
		opt(&res)
	}

	etag := generateETag(res)
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	w.Header().Set("ETag", etag)
	return writeJSON(w, res.Status, res)
}

func generateETag(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "error"
	}
	// Compute the SHA-256 hash of the JSON representation
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf(`"%x"`, hash)
}
