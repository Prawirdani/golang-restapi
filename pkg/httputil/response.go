package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"msg,omitempty"`
	Error   interface{} `json:"error,omitempty"`
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
	e := parseError(err)
	SendJson(w, e.Status, Response{Error: e.Cause})
}
