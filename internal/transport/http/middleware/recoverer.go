package middleware

import (
	"fmt"

	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

func PanicRecoverer(next handler.Func) handler.Func {
	return func(c *handler.Context) error {
		defer func() {
			if rec := recover(); rec != nil {
				err := fmt.Errorf("%v", rec)
				log.Error("panic recovered",
					err,
					"path", c.URLPath,
					"method", c.Method,
				)
				// httpx.HandleError(w, err)
			}
		}()

		return next(c)
	}
}
