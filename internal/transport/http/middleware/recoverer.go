package middleware

import (
	"fmt"
	"net/http"

	res "github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

// Recovery Middleware
func (c *Collection) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered",
					fmt.Errorf("%v", rec),
					"path", r.URL.Path,
					"method", r.Method,
				)
				res.HandleError(w, errors.InternalServer("internal server error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
