package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	apiErr "github.com/prawirdani/golang-restapi/pkg/errors"
)

var (
	ErrSessionExpired  = apiErr.Forbidden("session expired")
	ErrSessionNotFound = apiErr.Unauthorized("session not found, please login to proceed")
)

type Session struct {
	ID           int       `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	UserAgent    string    `db:"user_agent"`
	ExpiresAt    time.Time `db:"expires_at"`
	AccessedAt   time.Time `db:"accessed_at"`
}

func NewSession(userID uuid.UUID, userAgent string, expiry time.Duration) (Session, error) {
	if expiry <= 0 {
		return Session{}, errors.New("expiry must be greater than 0")
	}

	if userID == uuid.Nil {
		return Session{}, errors.New("user_id must not be empty")
	}

	// Generate a random 32 bytes for refresh token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return Session{}, err
	}

	refreshToken := hex.EncodeToString(bytes)

	currentTime := time.Now()
	sess := Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ExpiresAt:    currentTime.Add(expiry),
		AccessedAt:   currentTime,
	}

	return sess, nil
}

// IsExpired checks if the refresh token from the session has expired.
func (s Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
