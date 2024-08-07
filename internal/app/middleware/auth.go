package middleware

import (
	"net/http"
	"strings"

	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/response"
	"github.com/prawirdani/golang-restapi/pkg/token"
)

func (mw *Collection) AuthAccessToken(next http.Handler) http.Handler {
	return mw.authorize(common.AccessToken)(next)
}

func (mw *Collection) AuthRefreshToken(next http.Handler) http.Handler {
	return mw.authorize(common.RefreshToken)(next)
}

var ErrMissingToken = errors.Unauthorized("Missing auth token from cookie or Authorization bearer token")

// Token Authoriziation Middleware
func (mw *Collection) authorize(tt common.TokenType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the token string from the request cookie
			tokenStr := httputil.GetCookie(r, tt.Label())

			// If token doesn't exist in cookie, retrieve from Authorization header
			if tokenStr == "" {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = authHeader[len("Bearer "):]
				}
			}
			// If token is still empty, return an error
			if tokenStr == "" {
				response.HandleError(w, ErrMissingToken)
				return
			}

			payload, err := token.Decode(tokenStr, mw.cfg.Token.SecretKey)
			if err != nil {
				response.HandleError(w, err)
				return
			}

			// Check if the token type is the same as the expected token type
			if getTokenType(payload) != tt {
				response.HandleError(w, errors.Unauthorized("Invalid token type"))
				return
			}

			// Passing the map claims / payload to the next handler via Context.
			ctx := httputil.SetAuthCtx(r.Context(), payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTokenType(payload map[string]interface{}) common.TokenType {
	v := common.TokenType(payload["type"].(float64))
	return v
}
