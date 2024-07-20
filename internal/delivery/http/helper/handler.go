package helper

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

// Response writer for handling error
func HandleError(w http.ResponseWriter, err error) {
	e := errors.Parse(err)
	response := model.Response{
		Error: &model.ErrorBody{
			Code:    e.Status,
			Message: e.Message,
			Details: e.Cause,
		},
	}

	writeErr := httputil.WriteJSON(w, e.Status, response)
	if writeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
