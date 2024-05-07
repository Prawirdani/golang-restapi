package httputil

import "net/http"

type CustomHandlerFn func(w http.ResponseWriter, r *http.Request) error

// Custom Function handler wrapper to make it easier handling errors from api handler function.
func HandlerWrapper(fn CustomHandlerFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			HandleError(w, err)
			return
		}
	}
}

// Response writer for handling error
func HandleError(w http.ResponseWriter, err error) {
	e := parseErrors(err)
	response := HttpResponse{
		Error: &errorBody{
			Code:    e.status,
			Message: e.message,
			Details: e.cause,
		},
	}

	writeErr := writeJSON(w, e.status, response)
	if writeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
