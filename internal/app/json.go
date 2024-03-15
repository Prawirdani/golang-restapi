package app

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/model"
)

type JsonHandler struct{}

func NewJsonHandler() *JsonHandler {
	return &JsonHandler{}
}

// Send JSON Response
func (j *JsonHandler) Send(w http.ResponseWriter, status_code int, response model.Response) {
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
func (j *JsonHandler) Bind(r *http.Request, request interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(request); err != nil {
		return err
	}
	return nil
}
