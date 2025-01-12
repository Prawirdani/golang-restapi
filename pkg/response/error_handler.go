package response

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/errors"
)

// Response writer for handling error
func HandleError(w http.ResponseWriter, err error) {
	e := errors.Parse(err)
	response := ResponseBody{
		Message: e.Message,
		Errors:  e.Cause,
	}

	writeErr := writeJSON(w, e.Status, response)
	if writeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
