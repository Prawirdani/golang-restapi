package auth

import (
	"context"
)

type Repository interface {
	InsertSession(ctx context.Context, sess Session) error
	GetUserSessionBy(ctx context.Context, field string, searchVal any) (Session, error)
	DeleteSession(ctx context.Context, field string, val any) error
	DeleteExpiredSessions(ctx context.Context) error

	InsertResetPasswordToken(ctx context.Context, token ResetPasswordToken) error
	GetResetPasswordTokenObj(ctx context.Context, token string) (ResetPasswordToken, error)
	UseResetPasswordToken(ctx context.Context, tokenObj ResetPasswordToken) error
}
