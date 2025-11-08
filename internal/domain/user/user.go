package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/errorsx"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

const DEFAULT_PROFILE_IMG = "default-profile.jpg"

var (
	ErrEmailExist       = errorsx.Duplicate("email already exists")
	ErrUserNotFound     = errorsx.NotExists("user not found")
	ErrEmailNotVerified = errorsx.Forbidden("email is not registered or not verified")
)

type User struct {
	ID              uuid.UUID                  `db:"id"            json:"id"                validate:"required,uuid"`
	Name            string                     `db:"name"          json:"name"              validate:"required"`
	Email           string                     `db:"email"         json:"email"             validate:"required,email"`
	Phone           common.Nullable[string]    `db:"phone"         json:"phone"`
	Password        string                     `db:"password"      json:"-"                 validate:"required,min=8"`
	ProfileImage    string                     `db:"profile_image" json:"-"`
	CreatedAt       time.Time                  `db:"created_at"    json:"created_at"`
	UpdatedAt       time.Time                  `db:"updated_at"    json:"updated_at"`
	DeletedAt       common.Nullable[time.Time] `db:"deleted_at"    json:"-"`
	ProfileImageURL string                     `                   json:"profile_image_url"`
}

// Create new user from request payload
func New(i model.CreateUserInput, hashedPassword string) (*User, error) {
	u := User{
		ID:           uuid.New(),
		Name:         i.Name,
		Email:        i.Email,
		Phone:        common.NewNullable(i.Phone),
		ProfileImage: DEFAULT_PROFILE_IMG,
		Password:     hashedPassword,
	}

	if err := validator.Struct(u); err != nil {
		return nil, err
	}
	return &u, nil
}
