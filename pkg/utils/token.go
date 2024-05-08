package utils

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

var (
	ErrorTokenInvalid    = httputil.ErrUnauthorized("Invalid or expired token")
	ErrorTokenSignMethod = httputil.ErrUnauthorized("Invalid or mismatch token signing method")
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
)

type jwtClaims struct {
	Payload    map[string]interface{} `json:"payload"`
	payloadTag string                 `json:"-"`
	jwt.RegisteredClaims
}

// NewJwtClaims creates a new jwtClaims instance with the provided payload and payload tag name.
// The payload is a map of string keys and interface{} values for custom JSON data.
// The payloadTagName is the JSON tag used to serialize the payload in the JWT token.
func NewJwtClaims(payload map[string]interface{}, payloadTagName string) *jwtClaims {
	return &jwtClaims{
		Payload:    payload,
		payloadTag: payloadTagName,
	}
}

// MarshalJSON used for changing Payload json tag dynamicly,
// The MarshalJSON method is automatically called by the json.Marshal function when marshaling a struct to JSON.
// json.Marshal will be called by SignedString method.
func (jc *jwtClaims) MarshalJSON() ([]byte, error) {
	// encode the original
	m, _ := json.Marshal(*jc)

	// decode it back to get a map
	var a interface{}
	json.Unmarshal(m, &a)
	b := a.(map[string]interface{})

	// Replace the payload key
	b[jc.payloadTag] = b["payload"]
	delete(b, "payload")

	// Return encoding of the map
	return json.Marshal(b)
}

func GenerateToken(claims *jwtClaims, secret string, expiry time.Duration) (string, error) {
	timeNow := time.Now()

	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(timeNow)
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(timeNow.Add(expiry))

	// Sign JWT
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Parse, validate and returning the token map claims / payload.
func ParseToken(tokenString, secret string) (map[string]interface{}, error) {
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

	return *claims, err
}
