package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
)

// JWT Payload jwtClaims
type jwtClaims struct {
	UserID string `json:"userID,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, secret string, expiry time.Duration) (string, error) {
	timeNow := time.Now()
	claims := &jwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	// Sign JWT
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Parse and validate token and returning the token map claims / payload.
func ParseToken(tokenString, secret string) (map[string]interface{}, error) {
	claims := new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method != jwtSigningMethod {
			return nil, httputil.ErrUnauthorized("Invalid or mismatch token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, httputil.ErrUnauthorized("Invalid or expired token")
	}

	return *claims, err
}
