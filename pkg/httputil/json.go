package httputil

import (
	"encoding/json"
	"net/http"
)

// JSON Request body binder
func JsonBind(r *http.Request, data any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	return nil
}
