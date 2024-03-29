package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ID := uuid.New()
		tokenStr, err := GenerateToken(ID.String(), "secret-key", 5*time.Minute)
		require.Nil(t, err)
		require.NotEmpty(t, tokenStr)
	})
}

func TestParseToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ID := uuid.New()
		secretKey := "secret-key"

		tokenStr, err := GenerateToken(ID.String(), secretKey, 5*time.Minute)
		require.Nil(t, err)
		require.NotEmpty(t, tokenStr)

		mapClaims, err := ParseToken(tokenStr, secretKey)
		require.Nil(t, err)
		require.NotNil(t, mapClaims)

		require.Equal(t, ID.String(), mapClaims["userID"])
	})

	t.Run("expired-token", func(t *testing.T) {
		ID := uuid.New()
		secretKey := "secret-key"

		tokenStr, err := GenerateToken(ID.String(), secretKey, -5*time.Minute)
		require.Nil(t, err)
		require.NotEmpty(t, tokenStr)

		mapClaims, err := ParseToken(tokenStr, secretKey)
		require.NotNil(t, err)
		require.Nil(t, mapClaims)
		require.Equal(t, err, ErrorTokenInvalid)
	})
}
