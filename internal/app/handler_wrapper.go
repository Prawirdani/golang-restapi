package app

import (
	"net/http"

	respkg "github.com/prawirdani/golang-restapi/pkg/response"
)

// customHandler is a custom handler wrapper function type.
type customHandler func(w http.ResponseWriter, r *http.Request) error

// Custom Function handler wrapper to make it easier handling errors from handler function.
func handlerFn(fn customHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			respkg.HandleError(w, err)
			return
		}
	}
}
