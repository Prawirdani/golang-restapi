package middleware

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

// Token Authoriziation Middleware
func (mw *Collection) authorize(tt utils.TokenType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve, parse and validate the JWT token from the request.
			claims, err := utils.ParseJWT(r, &mw.cfg.Token, tt)
			if err != nil {
				httputil.HandleError(w, err)
				return
			}

			// Passing the map claims / payload to the next handler via Context.
			ctx := httputil.SetAuthCtx(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
