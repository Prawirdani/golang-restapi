package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	mockUserID := uuid.New()
	mockUserAgent := "user-agent"
	mockExpiry := 1 * time.Hour

	session, err := NewSession(mockUserID, mockUserAgent, mockExpiry)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, session.RefreshToken)
	assert.Equal(t, mockUserID, session.UserID)
	assert.Equal(t, mockUserAgent, session.UserAgent)
	assert.WithinDuration(t, time.Now().Add(mockExpiry), session.ExpiresAt, 1*time.Second)

	t.Run("Invalid-Expiry", func(t *testing.T) {
		_, err := NewSession(mockUserID, mockUserAgent, 0)
		require.Error(t, err)
		assert.Equal(t, "expiry must be greater than 0", err.Error())
	})

	t.Run("Invalid-UserID", func(t *testing.T) {
		_, err := NewSession(uuid.Nil, mockUserAgent, mockExpiry)
		require.Error(t, err)
		assert.Equal(t, "user_id must not be empty", err.Error())
	})

	t.Run("Expired", func(t *testing.T) {
		session, err := NewSession(mockUserID, mockUserAgent, mockExpiry)
		require.NoError(t, err)

		session.ExpiresAt = time.Now().Add(-1 * time.Hour)
		assert.True(t, session.IsExpired())
	})
}
