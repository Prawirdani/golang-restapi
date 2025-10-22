package response

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Utility function to help writing json to response body.
func writeJSON(w http.ResponseWriter, status int, response interface{}) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}
