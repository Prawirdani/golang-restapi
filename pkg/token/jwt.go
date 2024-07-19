package token

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

var (
	ErrorTokenInvalid    = httputil.ErrUnauthorized("Invalid or expired token")
	ErrorTokenSignMethod = httputil.ErrUnauthorized("Invalid or mismatch token signing method")
)

var (
	signingMethod = jwt.SigningMethodHS256
)

type Type uint8

const (
	Access Type = iota
	Refresh
)

// String returns the string representation of the token type and used as token cookie name.
func (t Type) String() string {
	return [...]string{"accessToken", "refreshToken"}[t]
}

type JWT struct {
	claims *jwtClaims
	value  string
	exp    time.Duration
	cookie *http.Cookie
}

// String returns the JWT token string.
func (j JWT) String() string {
	return j.value
}

// SetCookie sets the JWT token as a cookie in the response.
func (j JWT) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, j.cookie)
}

// Type returns the token type.
func (j JWT) Type() Type {
	return j.claims.TokenType
}

// TypeLabel returns the token type label.
func (j JWT) TypeLabel() string {
	return j.claims.TokenType.String()
}

// Payload returns the token payload.
func (j JWT) Payload() map[string]interface{} {
	return j.claims.User
}

type jwtClaims struct {
	User      map[string]interface{} `json:"user"`
	TokenType Type                   `json:"type"`
	jwt.RegisteredClaims
}
