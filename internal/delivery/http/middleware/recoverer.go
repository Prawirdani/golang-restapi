package middleware

import (
	"fmt"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/delivery/http/helper"
)

/* Panic recoverer middleware, it keep the service alive when crashes */
func (c *Collection) PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				helper.HandleError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
