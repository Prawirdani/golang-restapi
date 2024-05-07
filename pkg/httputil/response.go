package httputil

import (
	"net/http"
)

type HttpResponse struct {
	Data    any        `json:"data,omitempty"`
	Message *string    `json:"message,omitempty"`
	Status  int        `json:"-"`
	Error   *errorBody `json:"error,omitempty"`
}

type errorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func Response(w http.ResponseWriter, opts ...func(*HttpResponse)) error {
	res := &HttpResponse{
		Status: 200, // Default
	}

	for _, opt := range opts {
		opt(res)
	}

	return writeJSON(w, res.Status, res)
}
