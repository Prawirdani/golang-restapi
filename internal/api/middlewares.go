package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	chiCors "github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

// Main Router Middlewares

func cors(cfg *config.Config) func(next http.Handler) http.Handler {
	return chiCors.Handler(
		chiCors.Options{
			AllowedOrigins:   cfg.Cors.Origins,
			AllowCredentials: cfg.Cors.Credentials,
			Debug:            cfg.IsProduction(),
		},
	)
}

func gzip(next http.Handler) http.Handler {
	return chiMiddleware.Compress(6)(next)
}

func rateLimit(next http.Handler) http.Handler {
	return httprate.Limit(10, 10*time.Second, httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint))(next)
}

func reqLogger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}

/* Panic recoverer middleware, it keep the service alive when crashes */
func panicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				httputil.HandleError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
