package response

import "net/http"

type ResponseBody struct {
	Data    any        `json:"data,omitempty"`
	Message string     `json:"message,omitempty"`
	Status  int        `json:"-"`
	Error   *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

// Send is a function to send json response to client
// It uses option pattern to accepts multiple options to customize the response
func Send(w http.ResponseWriter, opts ...Option) error {
	res := ResponseBody{
		Status: http.StatusOK, // Default
	}
	for _, opt := range opts {
		opt(&res)
	}
	return writeJSON(w, res.Status, res)
}
