package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/spf13/viper"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
)

// JWT Payload jwtClaims
type jwtClaims struct {
	UserID string `json:"userID,omitempty"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTProvider(config *viper.Viper) *JWTProvider {
	return &JWTProvider{
		secretKey: []byte(config.GetString("jwt.secret_key")),
		expiry:    time.Duration(config.GetInt("jwt.expires") * int(time.Hour)),
	}
}

// Retrieve token string from request Auth Headers, parse it and return the claims/payload.
func (p *JWTProvider) VerifyRequest(r *http.Request) (map[string]interface{}, error) {
	// Retrieving from Request Auth Headers
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, httputil.ErrUnauthorized("Missing auth bearer token")
	}
	tokenString := authHeader[len("Bearer "):]

	claims, err := p.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Extract token claims/payload and check its validity
	return claims, nil
}

// Sign new token and return the token string.
func (p *JWTProvider) CreateToken(userID string) (*string, error) {
	timeNow := time.Now()
	claims := &jwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.expiry)),
		},
	}

	// Sign JWT
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	tokenStr, err := token.SignedString(p.secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}

// Parse and validate token and returning the token map claims / payload.
func (p *JWTProvider) ParseToken(tokenString string) (map[string]interface{}, error) {
	claims := new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method != jwtSigningMethod {
			return nil, httputil.ErrUnauthorized("Invalid or mismatch token signing method")
		}
		return p.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, httputil.ErrUnauthorized("Invalid or expired token")
	}

	return *claims, err
}
