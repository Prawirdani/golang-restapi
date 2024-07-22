package token

import (
	"time"

	"errors"

	"github.com/golang-jwt/jwt/v5"
	apiErr "github.com/prawirdani/golang-restapi/pkg/errors"
)

var (
	ErrTokenExpired = apiErr.Unauthorized("Authentication token has expired.")
	ErrInvalidToken = apiErr.Unauthorized("Invalid or malformed token.")
)

// claims represents the payload of a JWT token.
type claims struct {
	Payload map[string]interface{} `json:"payload"`
	jwt.RegisteredClaims
}

// Encode generates a new JWT token containing the given payload.
func Encode(secretKey string, payload map[string]interface{}, expiresIn time.Duration) (string, error) {
	currentTime := time.Now()
	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiresIn)),
		},
		Payload: make(map[string]interface{}),
	}
	for k, v := range payload {
		claims.Payload[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// Decode parses the given token string and returns the payload if the token is valid.
func Decode(tokenStr, secretKey string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	switch {
	case token != nil && token.Valid:
		break
	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, ErrTokenExpired
	default:
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	return claims.Payload, nil
}
