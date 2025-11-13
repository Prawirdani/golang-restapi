package user

import (
	"context"
)

// Repository defines persistence operations for user entities.
type Repository interface {
	// Store inserts a new user record.
	// Returns [ErrEmailExist] if the email already exists.
	Store(ctx context.Context, u *User) error

	// GetByID retrieves a user by their unique ID.
	// Returns [ErrNotFound] if no user exists with the given ID.
	GetByID(ctx context.Context, userID string) (*User, error)

	// GetByEmail retrieves a user by their email address.
	// Returns [ErrNotFound] if no user exists with the given email.
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update modifies an existing user record.
	// Returns [ErrEmailExist] if updating to an email that already exists.
	Update(ctx context.Context, u *User) error

	// Delete performs a soft delete by marking the user as deleted.
	Delete(ctx context.Context, u *User) error
}
