package httputil

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/prawirdani/golang-restapi/internal/model"
)

type CtxKey string

const AuthCtxKey CtxKey = "auth_ctx"

// SetAuthCtx sets the token payload to the context.
func SetAuthCtx(ctx context.Context, tokenPayload map[string]interface{}) context.Context {
	return context.WithValue(ctx, AuthCtxKey, tokenPayload)
}

// GetAuthCtx retrieves the token payload from the context.

type PayloadModel interface {
	model.AccessTokenPayload | model.RefreshTokenPayload
}

func GetAuthCtx[T PayloadModel](ctx context.Context) (T, error) {
	var payload T

	mapPayload, ok := ctx.Value(AuthCtxKey).(map[string]interface{})
	if !ok {
		return payload, errors.New("auth context not found")
	}

	jsonPayload, _ := json.Marshal(mapPayload)
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}
