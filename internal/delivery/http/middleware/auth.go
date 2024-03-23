package middleware

import (
	"context"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

type ClaimsKey string

const TOKEN_CLAIMS_CTX_KEY ClaimsKey = "token_claims"

// Token Authenticator Middleware
func (c *Collection) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.HandleError(w, httputil.ErrUnauthorized("Missing auth bearer token"))
			return
		}
		tokenString := authHeader[len("Bearer "):]

		claims, err := utils.ParseToken(tokenString, "secret")
		if err != nil {
			httputil.HandleError(w, err)
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
