package with

// With Package provide options for function that implement options pattern

import "github.com/prawirdani/golang-restapi/pkg/httputil"

type responseOption func(*httputil.HttpResponse)

func Data(data any) responseOption {
	return func(r *httputil.HttpResponse) {
		r.Data = data
	}
}
func Message(msg string) responseOption {
	return func(r *httputil.HttpResponse) {
		r.Message = &msg
	}
}

func Status(status int) responseOption {
	return func(r *httputil.HttpResponse) {
		r.Status = status
	}
}
