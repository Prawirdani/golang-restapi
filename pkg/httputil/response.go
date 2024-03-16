package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Send JSON HTTP Response
func SendJson(w http.ResponseWriter, status_code int, v any) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	_, _ = w.Write(buf.Bytes())
}

func SendError(w http.ResponseWriter, err error) {
	if e, ok := err.(*ApiError); ok {
		SendJson(w, e.Status, Response{Error: e.Error()})
		return
	}
	// Otherwise create an InternalServerError
	SendJson(w, http.StatusInternalServerError, Response{Error: err.Error()})
}
