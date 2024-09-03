package auth

import (
	"context"
	"encoding/json"
	"errors"
)

type CtxKey string

const AUTH_CTX_KEY CtxKey = "auth_ctx"

// SetContext sets the token payload to the context.
func SetContext(ctx context.Context, tokenPayload map[string]interface{}) context.Context {
	return context.WithValue(ctx, AUTH_CTX_KEY, tokenPayload)
}

// GetContext retrieves the token payload from the context.
func GetContext[T TokenPayload](ctx context.Context) (T, error) {
	var payload T

	mapPayload, ok := ctx.Value(AUTH_CTX_KEY).(map[string]interface{})
	if !ok {
		return payload, errors.New("auth context not found")
	}

	jsonPayload, _ := json.Marshal(mapPayload)
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}
