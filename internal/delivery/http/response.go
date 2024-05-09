package http

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type resOption func(*httputil.Response)

func data(v any) resOption {
	return func(r *httputil.Response) {
		r.Data = v
	}
}
func message(msg string) resOption {
	return func(r *httputil.Response) {
		r.Message = &msg
	}
}

func status(status int) resOption {
	return func(r *httputil.Response) {
		r.Status = status
	}
}

func Response(w http.ResponseWriter, opts ...func(*httputil.Response)) error {
	res := &httputil.Response{
		Status: 200, // Default
	}

	for _, opt := range opts {
		opt(res)
	}

	return httputil.WriteJSON(w, res.Status, res)
}
