package entity

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/stretchr/testify/require"
)

var registerPayload = model.RegisterRequest{
	Name:     "doe",
	Email:    "doe@mail.com",
	Password: "doe321",
}

func TestNewUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user, err := NewUser(registerPayload)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, registerPayload.Name, user.Name)
		require.Equal(t, registerPayload.Email, user.Email)
		require.NotEqual(t, registerPayload.Password, user.Password)
		require.NotEqual(t, uuid.Nil, user.ID)
	})

	t.Run("fail-missing-name", func(t *testing.T) {
		payload := registerPayload
		payload.Name = ""
		user, err := NewUser(payload)
		require.NotNil(t, err)
		require.Equal(t, User{}, user)
	})

	t.Run("fail-missing-email", func(t *testing.T) {
		payload := registerPayload
		payload.Email = ""

		user, err := NewUser(payload)
		require.NotNil(t, err)
		require.Equal(t, User{}, user)
	})

	t.Run("fail-missing-password", func(t *testing.T) {
		payload := registerPayload
		payload.Password = ""

		user, err := NewUser(payload)
		require.NotNil(t, err)
		require.Equal(t, User{}, user)

		_, ok := err.(validator.ValidationErrors)
		require.True(t, ok)
	})

	t.Run("fail-invalid-email", func(t *testing.T) {
		payload := registerPayload
		payload.Email = "invalid.email"

		user, err := NewUser(payload)
		require.NotNil(t, err)
		require.Equal(t, User{}, user)

		_, ok := err.(validator.ValidationErrors)
		require.True(t, ok)
	})

	t.Run("fail-minimum-password-chars", func(t *testing.T) {
		payload := registerPayload
		payload.Password = "123"

		user, err := NewUser(payload)
		require.NotNil(t, err)
		require.Equal(t, User{}, user)
		_, ok := err.(validator.ValidationErrors)
		require.True(t, ok)
	})
}

func TestVerifyPassword(t *testing.T) {
	user, err := NewUser(registerPayload)
	require.Nil(t, err)
	require.NotNil(t, user)

	t.Run("success", func(t *testing.T) {
		err := user.VerifyPassword(registerPayload.Password)
		require.Nil(t, err)
	})

	t.Run("wrong-password", func(t *testing.T) {
		err := user.VerifyPassword("wrong-pass")
		require.NotNil(t, err)
		require.Equal(t, err, ErrWrongCredentials)
	})
}
