package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateTokenType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tokenType := AccessToken
		err := tokenType.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		tokenType := TokenType(3)
		err := tokenType.Validate()
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTokenType, err)
	})
}

func TestGetTokenType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token, err := TokenEncode(SECRET, time.Minute*5, AccessToken, mockTokenPayload)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := TokenDecode(token, SECRET)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		tokenType, err := GetTokenType(payload)
		assert.NoError(t, err)
		assert.Equal(t, AccessToken, *tokenType)
	})

	t.Run("error", func(t *testing.T) {
		incompletePayload := map[string]interface{}{"user_id": "123"}
		tokenType, err := GetTokenType(incompletePayload)
		assert.Error(t, err)
		assert.Nil(t, tokenType)
		assert.Equal(t, ErrMissingTokenType, err)
	})
}
