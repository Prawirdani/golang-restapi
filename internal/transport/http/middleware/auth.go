package middleware

import (
	"strings"

	"github.com/prawirdani/golang-restapi/internal/domain/auth"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
)

func Auth(jwtSecret string) func(next handler.Func) handler.Func {
	return func(next handler.Func) handler.Func {
		return func(c *handler.Context) error {
			var tokenStr string

			// Try from cookie
			if cookie, err := c.GetCookie(handler.AccessTokenCookie); err == nil {
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
				return handler.ErrMissingAuthToken
			}

			// Validate token
			claims, err := auth.VerifyAccessToken(jwtSecret, tokenStr)
			if err != nil {
				return err
			}

			// Inject access token claims into request context
			ctx := auth.SetAccessTokenCtx(c.Context(), claims)
			c = c.WithContext(ctx)

			return next(c)
		}
	}
}
