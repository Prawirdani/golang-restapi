package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const SECRET = "secret"

var (
	accessTokenPayload = map[string]any{
		"user_id":  uuid.New().String(),
		"role":     "admin",
		"username": "doe",
	}
	refreshTokenPayload = map[string]any{"user_id": uuid.New().String()}
)

func TestTokenEncode(t *testing.T) {
	t.Run("AccessToken", func(t *testing.T) {
		token, err := TokenEncode(SECRET, time.Minute*5, AccessToken, accessTokenPayload)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		token, err := TokenEncode(SECRET, time.Hour*24*7, RefreshToken, refreshTokenPayload)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid-token-type", func(t *testing.T) {
		invalidTokenType := TokenType(99)
		token, err := TokenEncode(SECRET, time.Hour*24*7, invalidTokenType, refreshTokenPayload)
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestTokenDecode(t *testing.T) {
	t.Run("AccessToken", func(t *testing.T) {
		token, err := TokenEncode(SECRET, time.Minute*5, AccessToken, accessTokenPayload)
		require.NoError(t, err)

		payload, err := TokenDecode(token, SECRET)
		assert.NoError(t, err)
		assert.Equal(t, accessTokenPayload["user_id"], payload["user_id"])
		assert.Equal(t, accessTokenPayload["role"], payload["role"])
		assert.Equal(t, accessTokenPayload["username"], payload["username"])
		assert.Equal(t, AccessToken, TokenType(payload["typ"].(float64)))
	})

	t.Run("AccessToken-expired", func(t *testing.T) {
		token, err := TokenEncode(SECRET, -time.Minute*1, AccessToken, accessTokenPayload)
		require.NoError(t, err)

		_, err = TokenDecode(token, SECRET)
		assert.Error(t, err)
		assert.Equal(t, ErrTokenExpired, err)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		token, err := TokenEncode(SECRET, time.Hour*24*7, RefreshToken, refreshTokenPayload)
		require.NoError(t, err)

		payload, err := TokenDecode(token, SECRET)
		assert.NoError(t, err)
		assert.Equal(t, refreshTokenPayload["user_id"], payload["user_id"])
		assert.Equal(t, RefreshToken, TokenType(payload["typ"].(float64)))
	})

	t.Run("RefreshToken-expired", func(t *testing.T) {
		token, err := TokenEncode(SECRET, -time.Minute*1, RefreshToken, refreshTokenPayload)
		require.NoError(t, err)

		_, err = TokenDecode(token, SECRET)
		assert.Error(t, err)
		assert.Equal(t, ErrTokenExpired, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		payload, err := TokenDecode("invalid-token-string", SECRET)
		assert.Error(t, err)
		assert.Nil(t, payload)
	})
}
