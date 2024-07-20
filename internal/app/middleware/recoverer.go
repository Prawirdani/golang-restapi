package middleware

import (
	"fmt"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/response"
)

/* Panic recoverer middleware, it keep the service alive when crashes */
func (c *Collection) PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				response.HandleError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
