package middleware

import (
	"net/http"
	"strings"

	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/pkg/response"
)

// AccessToken authorization middleware
func (mw *Collection) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token string from the request cookie

		var tokenStr string

		if cookie, err := r.Cookie(auth.ACCESS_TOKEN); err == nil {
			tokenStr = cookie.Value
		}

		// If token doesn't exist in cookie, retrieve from Authorization header
		if tokenStr == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = authHeader[len("Bearer "):]
			}
		}

		// If token is still empty, return an error
		if tokenStr == "" {
			response.HandleError(w, auth.ErrMissingAccessToken)
			return
		}

		payload, err := auth.ValidateJWT(tokenStr, mw.cfg.Token.SecretKey)
		if err != nil {
			response.HandleError(w, err)
			return
		}

		// Passing the map claims / payload to the next handler via Context.
		ctx := auth.SetContext(r.Context(), payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
