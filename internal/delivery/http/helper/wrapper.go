package helper

import "net/http"

// CustomHandler is a custom handler wrapper function type.
type CustomHandler func(w http.ResponseWriter, r *http.Request) error

// Custom Function handler wrapper to make it easier handling errors from api handler function.
func HandlerFn(fn CustomHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			HandleError(w, err)
			return
		}
	}
}
