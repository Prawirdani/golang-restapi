package handler

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
	respkg "github.com/prawirdani/golang-restapi/pkg/response"
)

type responseOption func(*respkg.Base)

func data(v any) responseOption {
	return func(r *respkg.Base) {
		r.Data = v
	}
}

func message(msg string) responseOption {
	return func(r *respkg.Base) {
		r.Message = msg
	}
}

func status(status int) responseOption {
	return func(r *respkg.Base) {
		r.Status = status
	}
}

func response(w http.ResponseWriter, opts ...func(*respkg.Base)) error {
	res := &respkg.Base{
		Status: 200, // Default
	}

	for _, opt := range opts {
		opt(res)
	}

	return httputil.WriteJSON(w, res.Status, res)
}
