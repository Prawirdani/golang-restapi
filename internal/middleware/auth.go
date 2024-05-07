package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

type ClaimsKey string

const TOKEN_CLAIMS_CTX_KEY ClaimsKey = "token_claims"

// Token Authenticator Middleware
func (c *MiddlewareManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve Access Token from cookie header
		tokenString := httputil.GetCookie(r, c.tokenCfg.AccessCookieName)

		// If token doesn't exist in cookie, retrieve from Authorization header
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = authHeader[len("Bearer "):]
			}
		}

		// If token is still empty, return an error
		if tokenString == "" {
			httputil.HandleError(w, httputil.ErrUnauthorized("Missing auth token from cookie or Authorization bearer token"))
			return
		}

		claims, err := utils.ParseToken(tokenString, c.tokenCfg.SecretKey)
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
func (c *MiddlewareManager) GetAuthCtx(ctx context.Context) map[string]interface{} {
	return ctx.Value(TOKEN_CLAIMS_CTX_KEY).(map[string]interface{})
}
