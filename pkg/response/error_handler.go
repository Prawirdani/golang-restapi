package response

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

// Response writer for handling error
func HandleError(w http.ResponseWriter, err error) {
	e := errors.Parse(err)
	response := Base{
		Error: &ErrorBody{
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
