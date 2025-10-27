package response

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Utility function to help writing json to response body.
func writeJSON(w http.ResponseWriter, body *Body) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(body); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(body.Status)
	_, err := w.Write(buf.Bytes())
	return err
}
