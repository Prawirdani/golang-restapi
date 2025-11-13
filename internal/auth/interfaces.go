package auth

import (
	"context"

	"github.com/prawirdani/golang-restapi/internal/model"
)

// Repository Auth persistent storage interface
type Repository interface {
	// StoreSession saves a new session.
	StoreSession(ctx context.Context, sessionID *Session) error

	// GetSession returns a session by ID.
	GetSession(ctx context.Context, sessionID string) (*Session, error)

	// UpdateSession updates session data, mainly its expiration on logout.
	UpdateSession(ctx context.Context, session *Session) error

	// StoreResetPasswordToken saves a new password reset token.
	StoreResetPasswordToken(ctx context.Context, token *ResetPasswordToken) error

	// UpdateResetPasswordToken updates token data, mainly for update its 'used at' timestamp.
	UpdateResetPasswordToken(ctx context.Context, token *ResetPasswordToken) error

	// GetResetPasswordToken returns a token by its value.
	GetResetPasswordToken(ctx context.Context, tokenValue string) (*ResetPasswordToken, error)
}

// MessagePublisher Interface for message/event producer, eg rabbitmq, kafka etc.
type MessagePublisher interface {
	SendResetPasswordEmail(ctx context.Context, msg model.ResetPasswordEmailMessage) error
}
