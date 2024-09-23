package auth

import (
	"context"
	"errors"
)

var ErrTokenPayloadNotFound = errors.New("Access token payload not found in context")

type CtxKey string

const AUTH_CTX_KEY CtxKey = "auth_ctx"

// SetContext sets the access token payload to the context.
func SetContext(ctx context.Context, tokenPayload map[string]interface{}) context.Context {
	return context.WithValue(ctx, AUTH_CTX_KEY, tokenPayload)
}

// GetContext retrieves the access token payload from the context.
func GetContext(ctx context.Context) (map[string]interface{}, error) {
	payload, ok := ctx.Value(AUTH_CTX_KEY).(map[string]interface{})
	if !ok {
		return nil, ErrTokenPayloadNotFound
	}
	return payload, nil
}
