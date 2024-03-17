package middleware

import (
	"context"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

const TOKEN_CLAIMS_CTX_KEY = "token_claims"

// Token Authenticator Middleware
func (c *Collection) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := c.jwtAuth.VerifyRequest(r)

		if err != nil {
			httputil.SendError(w, err)
			return
		}

		// Passing the map claims / payload to the next handler via Context.
		ctx := context.WithValue(r.Context(), TOKEN_CLAIMS_CTX_KEY, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Util Function to retrieve Token map claims from the Request Context
func (c *Collection) GetAuthCtx(ctx context.Context) map[string]interface{} {
	return ctx.Value(TOKEN_CLAIMS_CTX_KEY).(map[string]interface{})
}
