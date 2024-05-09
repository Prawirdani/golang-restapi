package httputil

import "context"

const TOKEN_CLAIMS_CTX_KEY = "token_claims"

// Auth Context Setter
func SetAuthCtx(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, TOKEN_CLAIMS_CTX_KEY, claims)
}

// Retrieve Token map claims from the Request Context
func GetAuthCtx(ctx context.Context) map[string]interface{} {
	return ctx.Value(TOKEN_CLAIMS_CTX_KEY).(map[string]interface{})
}
