package middleware

import (
	"strings"

	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
	// httpx "github.com/prawirdani/golang-restapi/internal/transport/http"
)

// Auth AccessToken authorization middleware
//
//	func (mw *Collection) Auth(next http.Handler) http.Handler {
//		return handler.Middleware(mw.authHandler)(next)
//	}
func Auth(jwtSecret string) func(next handler.Func) handler.Func {
	return func(next handler.Func) handler.Func {
		return func(c *handler.Context) error {
			var tokenStr string

			// Try from cookie
			if cookie, err := c.GetCookie(auth.ACCESS_TOKEN); err == nil {
				tokenStr = cookie.Value
			}

			// Try from Authorization header
			if tokenStr == "" {
				authHeader := c.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = authHeader[len("Bearer "):]
				}
			}

			// If missing, return unauthorized error
			if tokenStr == "" {
				return auth.ErrTokenNotProvided
			}

			// Validate token
			payload, err := auth.ValidateJWT(tokenStr, jwtSecret)
			if err != nil {
				return err
			}

			// Inject auth payload into request context
			ctx := auth.SetAuthCtx(c.Context(), payload)
			c = c.WithContext(ctx)

			return next(c)
		}
	}
}
