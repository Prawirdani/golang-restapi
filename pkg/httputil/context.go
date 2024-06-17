package httputil

import "context"

type CtxKey string

const TOKEN_CLAIMS CtxKey = "token_claims"

// Auth Context Setter
func SetAuthCtx(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, TOKEN_CLAIMS, claims)
}

// Retrieve Token map claims from the Request Context
func GetAuthCtx(ctx context.Context) map[string]interface{} {
	return ctx.Value(TOKEN_CLAIMS).(map[string]interface{})
}
