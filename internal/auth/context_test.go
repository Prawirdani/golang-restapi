package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockTokenPayload = map[string]interface{}{
	"id":   "id-123",
	"name": "John Doe",
}

func TestAuthContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := SetContext(context.Background(), mockTokenPayload)
		payload, err := GetContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, mockTokenPayload, payload)
	})

	t.Run("ctx-not-exist", func(t *testing.T) {
		payload, err := GetContext(context.Background())
		assert.Error(t, err)
		assert.Nil(t, payload)
		assert.Equal(t, ErrTokenPayloadNotFound, err)
	})
}
