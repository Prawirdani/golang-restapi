// Package auth provides authentication and authorization functionality.
// This package handles user authentication through sessions, access tokens, and
// password management including secure hashing and password reset flows. It manages
// the complete authentication lifecycle from login through logout, including token
// generation, validation, and session management.
package auth

import "context"

// Repository defines the persistence operations for authentication data.
type Repository interface {
	// StoreSession creates a new session record.
	StoreSession(ctx context.Context, session *Session) error

	// GetSession retrieves a session by its ID.
	GetSession(ctx context.Context, sessionID string) (*Session, error)

	// UpdateSession updates an existing session (typically expiration).
	UpdateSession(ctx context.Context, session *Session) error

	// StoreResetPasswordToken creates a new password-reset token.
	StoreResetPasswordToken(ctx context.Context, token *ResetPasswordToken) error

	// UpdateResetPasswordToken updates an existing token (e.g., marking it used).
	UpdateResetPasswordToken(ctx context.Context, token *ResetPasswordToken) error

	// GetResetPasswordToken retrieves a token by its value.
	GetResetPasswordToken(ctx context.Context, value string) (*ResetPasswordToken, error)
}

// MessagePublisher defines the contract for publishing authentication-related
// messages to external systems (e.g., message queues, event buses).
// This enables asynchronous processing of notifications and events.
type MessagePublisher interface {
	// SendResetPasswordEmail publishes a message to trigger a password reset email.
	// The message contains the user's information and reset token that will be
	// consumed by an email service worker.
	// Returns an error if the message cannot be published to the queue.
	SendResetPasswordEmail(ctx context.Context, msg ResetPasswordEmailMessage) error
}
