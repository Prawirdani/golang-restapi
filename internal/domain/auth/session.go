// Package auth provides authentication and authorization functionality.
// This package handles user authentication through sessions, access tokens, and
// password management including secure hashing and password reset flows. It manages
// the complete authentication lifecycle from login through logout, including token
// generation, validation, and session management.
package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/domain"
)

var (
	ErrSessionExpired    = domain.ErrForbidden("Session expired")
	ErrSessionNotFound   = domain.ErrNotFound("Session not found, please login to proceed")
	ErrSessionEmptyUID   = errors.New("user_id must not be empty")
	ErrSessionInvalidTTL = errors.New("session ttl must be greater than 0")
)

// Session represents a long-lived refresh token stored server-side.
// The session ID is stored in an httpOnly cookie on the client and used to generate
// new access tokens when they expire. Sessions are stored in the database and can
// be revoked server-side for immediate logout.
//
// Architecture:
//   - Access Token (JWT): Short-lived (e.g., 15 min), sent with API requests
//   - Session ID: Long-lived (e.g., 30 days), stored in httpOnly cookie, used to refresh access tokens
//
// Security Features:
//   - Server-side revocation (logout invalidates session immediately)
//   - User agent tracking (detect token theft across devices)
//   - UUID v7 for session IDs (time-ordered for better DB performance)
type Session struct {
	// ID is the session identifier used as the refresh token.
	// This UUID v7 is exposed to clients via httpOnly cookies and used
	// to request new access tokens from refresh token endpoint.
	ID uuid.UUID `db:"id"`

	// UserID is the user who owns this session.
	UserID uuid.UUID `db:"user_id"`

	// UserAgent stores the client's User-Agent header for security tracking.
	// Used to detect session hijacking when requests come from different devices.
	UserAgent string `db:"user_agent"`

	// ExpiresAt is when this refresh token expires.
	// After expiration, the user must re-authenticate with credentials.
	ExpiresAt time.Time `db:"expires_at"`

	// AccessedAt tracks the last time this session was used to refresh an access token.
	// Updated on each successful refresh request for activity monitoring.
	AccessedAt time.Time `db:"accessed_at"`
}

// NewSession creates a new session for the given user.
// The session ID serves as a refresh token and should be stored in an httpOnly cookie.
func NewSession(
	userID uuid.UUID,
	userAgent string,
	ttl time.Duration,
) (*Session, error) {
	if ttl <= 0 {
		return nil, ErrSessionInvalidTTL
	}
	if userID == uuid.Nil {
		return nil, ErrSessionEmptyUID
	}

	sessID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sess := Session{
		ID:         sessID,
		UserID:     userID,
		UserAgent:  userAgent,
		ExpiresAt:  now.Add(ttl),
		AccessedAt: now,
	}
	return &sess, nil
}

// IsExpired checks if the session (refresh token) has expired.
// Expired sessions cannot be used to generate new access tokens and
// require the user to re-authenticate with their credentials.
//
// Returns true if the current time is past ExpiresAt.
func (s Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

// Revoke immediately expires the session by setting ExpiresAt to the current time.
func (s *Session) Revoke() {
	s.ExpiresAt = time.Now()
}
