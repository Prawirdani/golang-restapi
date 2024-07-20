package token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

// Parse the JWT token from the request and returns the claims. tokenType is used to determine the expected token type.
func ParseJWT(r *http.Request, cfg *config.TokenConfig, tokenType Type) (map[string]interface{}, error) {
	tokenString := httputil.GetCookie(r, tokenType.String())

	// If token doesn't exist in cookie, retrieve from Authorization header
	if tokenString == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[len("Bearer "):]
		}
	}

	// If token is still empty, return an error
	if tokenString == "" {
		return nil, errors.Unauthorized("Missing auth token from cookie or Authorization bearer token")
	}

	claims, err := parseJWT(tokenString, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	// Validate expected token type
	claimsTokenType, ok := claims["type"].(float64)
	if !ok {
		return nil, errors.Unauthorized("Invalid token type")
	}

	if Type(claimsTokenType) != tokenType {
		return nil, errors.Unauthorized(fmt.Sprintf("Invalid token type, expected %s", tokenType.String()))
	}

	return claims, nil
}

// Parse, validate and returning the token map claims / payload.
func parseJWT(tokenString, secret string) (map[string]interface{}, error) {
	claims := new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method != signingMethod {
			return nil, ErrorTokenSignMethod
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrorTokenInvalid
	}
	return *claims, nil
}
