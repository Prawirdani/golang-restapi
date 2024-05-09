package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// JSON Request body binder
func BindJSON(r *http.Request, data any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	return nil
}

// Utility function to help writing json to response body.
func WriteJSON(w http.ResponseWriter, status int, response interface{}) error {
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
