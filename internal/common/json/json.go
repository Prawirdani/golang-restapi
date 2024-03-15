package json

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/common/api"
)

// Send JSON HTTP Response
func Send(w http.ResponseWriter, status_code int, response api.Response) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	_, _ = w.Write(buf.Bytes())
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
