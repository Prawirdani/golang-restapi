package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	apiErr "github.com/prawirdani/golang-restapi/pkg/errors"
)

const (
	// Used as token cookie name and response body key.
	ACCESS_TOKEN  = "access_token"
	REFRESH_TOKEN = "refresh_token"
)

var (
	ErrTokenExpired = apiErr.Unauthorized("Token has expired.")
	ErrTokenInvalid = apiErr.Unauthorized("Invalid or malformed token.")
	ErrMissingToken = apiErr.Unauthorized(
		"Missing auth token from cookie or Authorization bearer token",
	)
	ErrEmptyTokenSecret = errors.New("secret key must not be empty")
)

// GenerateJWT generates a new Json Web Token containing the given payload. JWT is used as access token.
func GenerateJWT(
	secretKey string,
	expiry time.Duration,
	payload *map[string]any,
) (string, error) {
	if secretKey == "" {
		return "", ErrEmptyTokenSecret
	}

	currentTime := time.Now()

	mapClaims := jwt.MapClaims{
		"iat": jwt.NewNumericDate(currentTime),
		"exp": jwt.NewNumericDate(currentTime.Add(expiry)),
	}

	if payload != nil {
		for k, v := range *payload {
			mapClaims[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(secretKey))
}

// ValidateJWT decodes the given JWT string and returns the map payload.
func ValidateJWT(tokenStr, secretKey string) (map[string]any, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apiErr.BadRequest(
				fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]),
			)
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
