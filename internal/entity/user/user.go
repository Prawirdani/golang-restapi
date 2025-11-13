package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/errorsx"
	"github.com/prawirdani/golang-restapi/pkg/nullable"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

const DefaultProfilePictureFile = "default-profile.jpg"

var (
	ErrEmailExist       = errorsx.Duplicate("email already exists")
	ErrNotFound         = errorsx.NotExists("user not found")
	ErrEmailNotVerified = errorsx.Forbidden("email is not registered or not verified")
)

type User struct {
	ID              uuid.UUID                    `db:"id"            json:"id"                validate:"required,uuid"`
	Name            string                       `db:"name"          json:"name"              validate:"required"`
	Email           string                       `db:"email"         json:"email"             validate:"required,email"`
	Phone           nullable.Nullable[string]    `db:"phone"         json:"phone"`
	Password        string                       `db:"password"      json:"-"                 validate:"required,min=8"`
	ProfileImage    string                       `db:"profile_image" json:"-"`
	CreatedAt       time.Time                    `db:"created_at"    json:"created_at"`
	UpdatedAt       time.Time                    `db:"updated_at"    json:"updated_at"`
	DeletedAt       nullable.Nullable[time.Time] `db:"deleted_at"    json:"-"`
	ProfileImageURL string                       `                   json:"profile_image_url"`
}

// New Create new user
func New(inp model.CreateUserInput, hashedPassword string) (*User, error) {
	u := User{
		ID:           uuid.New(),
		Name:         inp.Name,
		Email:        inp.Email,
		Phone:        nullable.New(inp.Phone),
		ProfileImage: DefaultProfilePictureFile,
		Password:     hashedPassword,
	}

	if err := validator.Struct(u); err != nil {
		return nil, err
	}
	return &u, nil
}
