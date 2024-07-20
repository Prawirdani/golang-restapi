package handler

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type responseOption func(*model.Response)

func data(v any) responseOption {
	return func(r *model.Response) {
		r.Data = v
	}
}

func message(msg string) responseOption {
	return func(r *model.Response) {
		r.Message = msg
	}
}

func status(status int) responseOption {
	return func(r *model.Response) {
		r.Status = status
	}
}

func response(w http.ResponseWriter, opts ...func(*model.Response)) error {
	res := &model.Response{
		Status: 200, // Default
	}

	for _, opt := range opts {
		opt(res)
	}

	return httputil.WriteJSON(w, res.Status, res)
}
