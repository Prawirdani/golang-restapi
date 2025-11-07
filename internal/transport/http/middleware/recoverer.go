package middleware

import (
	"fmt"
	"net/http"

	httpx "github.com/prawirdani/golang-restapi/internal/transport/http"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

// Recovery Middleware
func (c *Collection) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err := fmt.Errorf("%w", rec)
				log.Error("panic recovered",
					err,
					"path", r.URL.Path,
					"method", r.Method,
				)
				httpx.HandleError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
