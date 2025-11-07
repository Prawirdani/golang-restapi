package response

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonBody struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	status  int    `json:"-"`
}

// JSON is a function to send json response to client
// It uses option pattern to accepts multiple options to customize the response
func JSON(w http.ResponseWriter, r *http.Request, opts ...Option) error {
	res := jsonBody{
		status: http.StatusOK, // Default
	}
	for _, opt := range opts {
		opt(&res)
	}

	// Only use ETag for successful responses (2xx)
	if res.status >= 200 && res.status < 300 {
		etag := eTag(res)
		if etag != "" {
			// Check If-None-Match header
			if match := r.Header.Get("If-None-Match"); match == etag {
				w.WriteHeader(http.StatusNotModified)
				return nil
			}
			w.Header().Set("ETag", etag)
			w.Header().Set("Cache-Control", "private, must-revalidate")
		}
		w.Header().Set("ETag", etag)
	}

	return WriteJson(w, res.status, &res)
}

func eTag(data any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	h := sha256.Sum256(b)
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(h[:]))
}

// Utility function to help writing json to response body.
func WriteJson(w http.ResponseWriter, status int, body any) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}
