// Package auth provides authentication and authorization functionality.
// This package handles user authentication through sessions, access tokens, and
// password management including secure hashing and password reset flows. It manages
// the complete authentication lifecycle from login through logout, including token
// generation, validation, and session management.
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/golang-restapi/internal/domain"
)

var (
	ErrAccessTokenExpired        = domain.ErrUnauthorized("Access token expired")
	ErrAccessTokenClaimsNotFound = errors.New("access token claims not found in context")
)

type AccessTokenClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

// SignAccessToken generates a new JWT for access token
func SignAccessToken(
	secretKey string,
	claims AccessTokenClaims,
	ttl time.Duration,
) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// VerifyAccessToken parses and validates the token, returning the claims if valid.
func VerifyAccessToken(secretKey, tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AccessTokenClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrAccessTokenExpired
		}
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, errors.New("invalid token claims type")
	}

	return claims, nil
}

type accessTokenCtxKey struct{}

var atCtx accessTokenCtxKey

// SetAccessTokenCtx sets the access token jwt claims to the context.
func SetAccessTokenCtx(ctx context.Context, claims *AccessTokenClaims) context.Context {
	return context.WithValue(ctx, atCtx, claims)
}

// GetAccessTokenCtx retrieves the access token jwt claims from the context.
func GetAccessTokenCtx(ctx context.Context) (*AccessTokenClaims, error) {
	claims, ok := ctx.Value(atCtx).(*AccessTokenClaims)
	if !ok {
		return nil, ErrAccessTokenClaimsNotFound
	}
	return claims, nil
}
