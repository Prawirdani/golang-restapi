package middleware

import (
	"fmt"
	"net/http"

	res "github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

/* Panic recoverer middleware, it keep the service alive when crashes */
func (c *Collection) PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				c.logger.Error(
					logging.RuntimePanic,
					"middleware.PanicRecoverer",
					fmt.Sprintf("%v", rvr),
				)
				res.HandleError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
