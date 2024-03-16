package json

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Send JSON HTTP Response
func Send(w http.ResponseWriter, status_code int, v any) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	_, err := w.Write(buf.Bytes())

	return err
}

// JSON Request body binder
func Bind(r *http.Request, data any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	return nil
}
