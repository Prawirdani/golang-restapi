package user

import (
	"context"
)

type Repository interface {
	Insert(ctx context.Context, u User) error
	// GetUsers(ctx context.Context) ([]User, error)
	GetUserBy(ctx context.Context, field string, value any) (User, error)
	UpdateUser(ctx context.Context, u User) error
	DeleteUser(ctx context.Context, userId string) error
}
