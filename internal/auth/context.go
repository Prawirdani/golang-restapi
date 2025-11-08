package auth

import (
	"context"
	"errors"
)

var ErrTokenPayloadNotFound = errors.New("access token payload not found in context")

type ctxKey struct{}

var authCtxKey ctxKey

// SetAuthCtx sets the access token payload to the context.
func SetAuthCtx(ctx context.Context, tokenPayload map[string]any) context.Context {
	return context.WithValue(ctx, authCtxKey, tokenPayload)
}

// GetAuthCtx retrieves the access token payload from the context.
func GetAuthCtx(ctx context.Context) (map[string]any, error) {
	payload, ok := ctx.Value(authCtxKey).(map[string]any)
	if !ok {
		return nil, ErrTokenPayloadNotFound
	}
	return payload, nil
}
