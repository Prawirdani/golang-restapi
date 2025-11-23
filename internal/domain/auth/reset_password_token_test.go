package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResetPasswordToken(t *testing.T) {
	userID := uuid.New()
	ttl := time.Hour

	token, err := NewResetPasswordToken(userID, ttl)
	require.NoError(t, err)
	assert.NotEmpty(t, token.Value)
	assert.Equal(t, userID, token.UserID)
	assert.False(t, token.Expired())
	assert.False(t, token.Used())
	assert.True(t, token.ExpiresAt.After(time.Now()))
}

func TestResetPasswordTokenExpired(t *testing.T) {
	userID := uuid.New()

	t.Run("not expired", func(t *testing.T) {
		token, err := NewResetPasswordToken(userID, time.Hour)
		require.NoError(t, err)
		assert.False(t, token.Expired())
	})

	t.Run("expired", func(t *testing.T) {
		token, err := NewResetPasswordToken(userID, -time.Hour)
		require.NoError(t, err)
		assert.True(t, token.Expired())
	})
}

func TestResetPasswordTokenUsed(t *testing.T) {
	userID := uuid.New()

	t.Run("not used", func(t *testing.T) {
		token, err := NewResetPasswordToken(userID, time.Hour)
		require.NoError(t, err)
		assert.False(t, token.Used())
	})

	t.Run("used", func(t *testing.T) {
		token, err := NewResetPasswordToken(userID, time.Hour)
		require.NoError(t, err)
		token.Revoke()
		assert.True(t, token.Used())
	})
}

func TestResetPasswordTokenRevoke(t *testing.T) {
	userID := uuid.New()
	token, err := NewResetPasswordToken(userID, time.Hour)
	require.NoError(t, err)

	assert.False(t, token.Used())
	token.Revoke()
	assert.True(t, token.Used())
	assert.True(t, token.UsedAt.NotNull())
}
