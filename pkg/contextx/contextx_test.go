package contextx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockPayload = map[string]any{
	"id":   "id-123",
	"name": "John Doe",
}

func TestAuthContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := SetAuthCtx(context.Background(), mockPayload)
		payload, err := GetAuthCtx(ctx)
		assert.NoError(t, err)
		assert.Equal(t, mockPayload, payload)
	})

	t.Run("ctx-not-exist", func(t *testing.T) {
		payload, err := GetAuthCtx(context.Background())
		assert.Error(t, err)
		assert.Nil(t, payload)
		assert.Equal(t, ErrTokenPayloadNotFound, err)
	})
}
