package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/pkg/errorsx"
	"github.com/prawirdani/golang-restapi/pkg/nullable"
)

var (
	// ErrResetPasswordTokenInvalid is returned when the reset token is invalid or expired.
	ErrResetPasswordTokenInvalid = errorsx.Forbidden(
		"The reset token is invalid or expired. Please request a new password reset.",
	)

	// ErrResetPasswordTokenNotFound is returned when no matching reset token exists.
	ErrResetPasswordTokenNotFound = errorsx.NotExists("Reset password token not found")
)

// ResetPasswordToken represents a one-time password reset token tied to a user.
type ResetPasswordToken struct {
	UserID    uuid.UUID                    `db:"user_id"    json:"user_id"`
	Value     string                       `db:"value"      json:"value"`
	ExpiresAt time.Time                    `db:"expires_at" json:"expires_at"`
	UsedAt    nullable.Nullable[time.Time] `db:"used_at"    json:"used_at"`
}

// NewResetPasswordToken creates a new token for the given user with a specified expiration.
func NewResetPasswordToken(userID uuid.UUID, ttl time.Duration) (*ResetPasswordToken, error) {
	bs := make([]byte, 32)
	if _, err := rand.Read(bs); err != nil {
		return nil, err
	}

	return &ResetPasswordToken{
		UserID:    userID,
		Value:     hex.EncodeToString(bs),
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

// Expired reports whether the token has passed its expiration time.
func (t ResetPasswordToken) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// Used reports whether the token has already been used.
func (t ResetPasswordToken) Used() bool {
	return t.UsedAt.NotNull()
}

// Revoke marks the token as used immediately.
func (t *ResetPasswordToken) Revoke() {
	t.UsedAt = nullable.New(time.Now())
}
