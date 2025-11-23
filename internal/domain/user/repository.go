// Package user provides the domain model and business logic for managing users in system.
package user

import (
	"context"
)

// Repository defines the contract for user data persistence operations.
type Repository interface {
	// Store creates a new user record.
	// Returns [ErrEmailExists] if a user with the same email already exists.
	Store(ctx context.Context, u *User) error

	// GetByID retrieves a user by their unique identifier.
	// Returns [ErrNotFound] if no user exists with the given ID.
	GetByID(ctx context.Context, userID string) (*User, error)

	// GetByEmail retrieves a user by their email address.
	// Returns [ErrNotFound] if no user exists with the given email.
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update modifies an existing user record.
	// Returns [ErrEmailExists] if updating to an email that already exists.
	Update(ctx context.Context, u *User) error
}
