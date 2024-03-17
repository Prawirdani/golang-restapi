package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/model"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) ToResponse() *model.UserResponse {
	return &model.UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
