package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const SECRET = "secret"

var accessTokenPayload = map[string]any{
	"user_id":  "id-123",
	"username": "doe",
}

func TestTokenEncode(t *testing.T) {
	token, err := GenerateJWT(SECRET, time.Minute*5, &accessTokenPayload)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	t.Run("EmptySecret", func(t *testing.T) {
		token, err := GenerateJWT("", time.Minute*1, &accessTokenPayload)
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, ErrEmptyTokenSecret, err)
	})
}

func TestTokenDecode(t *testing.T) {
	token, err := GenerateJWT(SECRET, time.Minute*5, &accessTokenPayload)
	require.NoError(t, err)

	payload, err := ValidateJWT(token, SECRET)
	assert.NoError(t, err)
	assert.Equal(t, accessTokenPayload["user_id"], payload["user_id"])
	assert.Equal(t, accessTokenPayload["username"], payload["username"])

	t.Run("Expired", func(t *testing.T) {
		token, err := GenerateJWT(SECRET, -time.Minute*1, &accessTokenPayload)
		require.NoError(t, err)

		_, err = ValidateJWT(token, SECRET)
		assert.Error(t, err)
		assert.Equal(t, ErrTokenExpired, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		payload, err := ValidateJWT("invalid-token-string", SECRET)
		assert.Error(t, err)
		assert.Nil(t, payload)
	})
}
