package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/spf13/viper"
)

var (
	jwtSigningMethod          = jwt.SigningMethodHS256
	missingAuthHeaderError    = errors.New("Missing auth bearer token")
	invalidSigningMethodError = errors.New("Invalid signing method")
	invalidTokenError         = errors.New("Invalid or expired token")
)

// JWT Payload Claims
type Claims struct {
	UserID string `json:"userID,omitempty"`
	jwt.RegisteredClaims
}

type Provider struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTProvider(config *viper.Viper) *Provider {
	return &Provider{
		secretKey: []byte(config.GetString("jwt.secret_key")),
		expiry:    time.Duration(config.GetInt("jwt.expires") * int(time.Hour)),
	}
}

// Retrieve token string from request Auth Headers, parse it and return the claims/payload.
func (p *Provider) ValidateFromRequest(r *http.Request) (map[string]interface{}, error) {
	// Retrieving from Request Auth Headers
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, httputil.ErrUnauthorized(missingAuthHeaderError)
	}
	tokenString := authHeader[len("Bearer "):]

	claims, err := p.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Extract token claims/payload and check its validity
	return claims, nil
}

// Sign new token and return the token string.
func (p *Provider) CreateToken(u *entity.User) (*string, error) {
	timeNow := time.Now()
	claims := &Claims{
		UserID: u.ID.String(),
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
func (p *Provider) parseToken(tokenString string) (map[string]interface{}, error) {
	claims := new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidSigningMethodError
		} else if method != jwtSigningMethod {
			return nil, invalidSigningMethodError
		}
		return p.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, invalidTokenError
	}

	return *claims, err
}
