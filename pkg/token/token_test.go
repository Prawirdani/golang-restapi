package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const secretKey = "secret"

var payload = map[string]any{"user_id": 1, "role": "admin", "username": "doe"}

func TestEncode(t *testing.T) {
	expiresIn := time.Minute * 5
	token, err := Encode(secretKey, payload, expiresIn)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestDecode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expiresIn := time.Minute * 5
		token, err := Encode(secretKey, payload, expiresIn)
		assert.Nil(t, err)

		decodedPayload, err := Decode(token, secretKey)
		assert.Nil(t, err)

		// Convert any int type because default behavior of encoding/json when unmarshalling json numbers into interface{} will be converted in to float64.
		assert.Equal(t, float64(payload["user_id"].(int)), decodedPayload["user_id"])
		assert.Equal(t, payload["role"], decodedPayload["role"])
		assert.Equal(t, payload["username"], decodedPayload["username"])
	})

	t.Run("expired", func(t *testing.T) {
		expiresIn := -time.Minute * 1
		token, err := Encode(secretKey, payload, expiresIn)
		assert.Nil(t, err)

		_, err = Decode(token, secretKey)
		assert.Equal(t, ErrTokenExpired, err)
	})

	t.Run("invalid", func(t *testing.T) {
		token := "invalid-token-string"
		_, err := Decode(token, secretKey)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInvalidToken, err)
	})
}
