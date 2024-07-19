package middleware

import (
	"fmt"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

/* Panic recoverer middleware, it keep the service alive when crashes */
func (m *Collection) PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				httputil.HandleError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
