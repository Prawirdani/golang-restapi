package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/pkg/common"
)

var (
	ErrResetPasswordTokenInvalid = errors.New(
		"The reset token is invalid or expired. Please request a new password reset.",
	)
	ErrResetPasswordTokenNotFound = errors.New(
		"Reset password token not found",
	)
)

type ResetPasswordToken struct {
	UserId    uuid.UUID                  `db:"user_id"    json:"user_id"`
	Value     string                     `db:"value"      json:"value"`
	ExpiresAt time.Time                  `db:"expires_at" json:"expires_at"`
	UsedAt    common.Nullable[time.Time] `db:"used_at"    json:"used_at"`
}

func (t ResetPasswordToken) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func (t ResetPasswordToken) Used() bool {
	// Valid means not null, means its been used
	return t.UsedAt.Valid
}

func NewResetPasswordToken(userId uuid.UUID, expiresAt time.Time) (*ResetPasswordToken, error) {
	bs := make([]byte, 32)
	if _, err := rand.Read(bs); err != nil {
		return nil, err
	}
	value := hex.EncodeToString(bs)

	token := &ResetPasswordToken{
		UserId:    userId,
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return token, nil
}
