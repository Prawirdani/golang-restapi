package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := NewUser(model.RegisterRequest{
			Name:     "doe",
			Email:    "doe@mail.com",
			Password: "doe321",
		})

		require.NotNil(t, user)
		require.Equal(t, "doe", user.Name)
		require.Equal(t, "doe@mail.com", user.Email)
		require.Equal(t, "doe321", user.Password)
		require.NotEqual(t, uuid.Nil, user.ID)
	})

}

func TestValidate(t *testing.T) {
	user := NewUser(model.RegisterRequest{
		Name:     "doe",
		Email:    "doe@mail.com",
		Password: "doe321",
	})
	t.Run("success", func(t *testing.T) {
		newUser := user
		err := newUser.Validate()
		require.Nil(t, err)
	})

	t.Run("fail-missing-name", func(t *testing.T) {
		newUser := user
		newUser.Name = ""
		err := newUser.Validate()
		require.NotNil(t, err)
	})
	t.Run("fail-missing-email", func(t *testing.T) {
		newUser := user
		newUser.Email = ""
		err := newUser.Validate()
		require.NotNil(t, err)
	})
	t.Run("fail-missing-password", func(t *testing.T) {
		newUser := user
		newUser.Password = ""
		err := newUser.Validate()
		require.NotNil(t, err)
	})
	t.Run("fail-invalid-email", func(t *testing.T) {
		newUser := user
		newUser.Email = "doe.mail.com"
		err := newUser.Validate()
		require.NotNil(t, err)
	})
	t.Run("fail-minimum-password-chars", func(t *testing.T) {
		newUser := user
		newUser.Password = "12345"
		err := newUser.Validate()
		require.NotNil(t, err)
	})
}

func TestEncryptPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		newUser := NewUser(model.RegisterRequest{
			Name:     "john doe",
			Email:    "doe@mail.com",
			Password: "doe321",
		})

		err := newUser.Validate()
		require.Nil(t, err)

		err = newUser.EncryptPassword()
		require.Nil(t, err)
	})
}

func TestVerifyPassword(t *testing.T) {
	user := NewUser(model.RegisterRequest{
		Name:     "john doe",
		Email:    "doe@mail.com",
		Password: "doe321",
	})

	err := user.Validate()
	require.Nil(t, err)

	err = user.EncryptPassword()
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		err := user.VerifyPassword("doe321")
		require.Nil(t, err)
	})

	t.Run("wrong-password", func(t *testing.T) {
		err := user.VerifyPassword("wrong-pass")
		require.NotNil(t, err)
		require.Equal(t, err, ErrorWrongCredentials)
	})
}

func TestGenerateToken(t *testing.T) {
	user := NewUser(model.RegisterRequest{
		Name:     "doe",
		Email:    "doe@mail.com",
		Password: "doe321",
	})
	require.NotNil(t, user)

	err := user.Validate()
	require.Nil(t, err)

	err = user.EncryptPassword()
	require.Nil(t, err)

	tokenStr, err := user.GenerateToken("secret-key")
	require.Nil(t, err)
	require.NotEmpty(t, tokenStr)
}
