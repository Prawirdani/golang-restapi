package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	apiErr "github.com/prawirdani/golang-restapi/pkg/errors"
)

var (
	ErrTokenExpired = apiErr.Unauthorized("Token has expired.")
	ErrInvalidToken = apiErr.Unauthorized("Invalid or malformed token.")
)

// TokenEncode generates a new JWT token containing the given payload.
func TokenEncode(
	secretKey string,
	expiry time.Duration,
	tokenType TokenType,
	payload map[string]interface{},
) (string, error) {
	if err := tokenType.Validate(); err != nil {
		return "", err
	}

	currentTime := time.Now()

	mapClaims := jwt.MapClaims{
		"iat": jwt.NewNumericDate(currentTime),
		"exp": jwt.NewNumericDate(currentTime.Add(expiry)),
		"typ": tokenType,
	}

	for k, v := range payload {
		mapClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(secretKey))
}

// TokenDecode decodes the given JWT token string and returns the map payload.
func TokenDecode(tokenStr, secretKey string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
