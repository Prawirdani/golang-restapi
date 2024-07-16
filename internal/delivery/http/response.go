package http

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type responseOption func(*httputil.Response)

func data(v any) responseOption {
	return func(r *httputil.Response) {
		r.Data = v
	}
}

func message(msg string) responseOption {
	return func(r *httputil.Response) {
		r.Message = msg
	}
}

func status(status int) responseOption {
	return func(r *httputil.Response) {
		r.Status = status
	}
}

func response(w http.ResponseWriter, opts ...func(*httputil.Response)) error {
	res := &httputil.Response{
		Status: 200, // Default
	}

	for _, opt := range opts {
		opt(res)
	}

	return httputil.WriteJSON(w, res.Status, res)
}
