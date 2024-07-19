package token

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/config"
)

func GenerateJWT(cfg *config.Config, payload map[string]interface{}, tokenType Type) (JWT, error) {
	timeNow := time.Now()

	var expiry time.Duration
	if tokenType == Access {
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
	token := jwt.NewWithClaims(signingMethod, claims)
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
