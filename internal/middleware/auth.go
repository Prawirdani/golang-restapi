package middleware

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

// Token Authenticator Middleware
func (c *MiddlewareManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve, Parse, Validate access token from request
		claims, err := utils.ParseJWT(r, &c.cfg.Token, utils.AccessToken)
		if err != nil {
			httputil.HandleError(w, err)
			return
		}

		// Passing the map claims / payload to the next handler via Context.
		ctx := httputil.SetAuthCtx(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
