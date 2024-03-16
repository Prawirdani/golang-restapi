package httputil

import (
	"net/http"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func HandlerWrapper(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			// If error was the object of ApiError
			if e, ok := err.(*ApiError); ok {
				JsonSend(w, e.Status, Response{Error: e.Error()})
				return
			}
			// Otherwise create an InternalServerError
			JsonSend(w, http.StatusInternalServerError, Response{Error: err.Error()})
		}

	}
}
