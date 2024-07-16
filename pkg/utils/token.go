package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

var (
	ErrorTokenInvalid    = httputil.ErrUnauthorized("Invalid or expired token")
	ErrorTokenSignMethod = httputil.ErrUnauthorized("Invalid or mismatch token signing method")
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
)

type TokenType uint8

const (
	AccessToken TokenType = iota
	RefreshToken
)

// String returns the string representation of the token type and used as token cookie name.
func (t TokenType) String() string {
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
func (j JWT) Type() TokenType {
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
	TokenType TokenType              `json:"type"`
	jwt.RegisteredClaims
}

func GenerateJWT(cfg *config.Config, payload map[string]interface{}, tokenType TokenType) (JWT, error) {
	timeNow := time.Now()

	var expiry time.Duration
	if tokenType == AccessToken {
		expiry = time.Duration(cfg.Token.AccessTokenExpiry) * time.Minute
	} else {
		expiry = time.Duration(cfg.Token.RefreshTokenExpiry) * time.Hour * 24
	}

	claims := &jwtClaims{
		User:      payload,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(expiry)),
		},
	}

	// Sign JWT
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	tokenStr, err := token.SignedString([]byte(cfg.Token.SecretKey))
	if err != nil {
		return JWT{}, nil
	}

	// Create cookie
	cookie := &http.Cookie{
		Name:     tokenType.String(),
		Value:    tokenStr,
		Path:     "/",
		Expires:  timeNow.Add(expiry),
		HttpOnly: cfg.IsProduction(),
	}

	return JWT{
		value:  tokenStr,
		claims: claims,
		exp:    expiry,
		cookie: cookie,
	}, nil
}

// Parse the JWT token from the request and returns the claims. tokenType is used to determine the expected token type.
func ParseJWT(r *http.Request, cfg *config.TokenConfig, tokenType TokenType) (map[string]interface{}, error) {
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
		return nil, httputil.ErrUnauthorized("Missing auth token from cookie or Authorization bearer token")
	}

	claims, err := parseJWT(tokenString, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	// Validate expected token type
	claimsTokenType, ok := claims["type"].(float64)
	if !ok {
		return nil, httputil.ErrUnauthorized("Invalid token type")
	}

	if TokenType(claimsTokenType) != tokenType {
		return nil, httputil.ErrUnauthorized(fmt.Sprintf("Invalid token type, expected %s", tokenType.String()))
	}

	return claims, nil
}

// Parse, validate and returning the token map claims / payload.
func parseJWT(tokenString, secret string) (map[string]interface{}, error) {
	claims := new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method != jwtSigningMethod {
			return nil, ErrorTokenSignMethod
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrorTokenInvalid
	}
	return *claims, nil
}
