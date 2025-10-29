package auth

import (
	"context"
)

type Repository interface {
	InsertSession(ctx context.Context, sess *Session) error
	GetUserSessionBy(ctx context.Context, field string, value any) (*Session, error)
	DeleteSession(ctx context.Context, field string, value any) error
	DeleteExpiredSessions(ctx context.Context) error

	InsertResetPasswordToken(ctx context.Context, token *ResetPasswordToken) error
	GetResetPasswordTokenObj(ctx context.Context, tokenValue string) (*ResetPasswordToken, error)
	InvalidateResetPasswordToken(ctx context.Context, token *ResetPasswordToken) error
}
