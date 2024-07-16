package entity

import (
	"log"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/utils"
	"github.com/stretchr/testify/require"
)

var cfg *config.Config

func init() {
	config, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatal(err)
	}
	cfg = config
}

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
		require.Equal(t, err, ErrorWrongCredentials)
	})
}

func TestGenerateToken(t *testing.T) {
	user, err := NewUser(registerPayload)
	require.Nil(t, err)
	require.NotNil(t, user)

	t.Run("AccessToken", func(t *testing.T) {
		token, err := user.GenerateAccessToken(cfg)
		require.Nil(t, err)
		require.NotEmpty(t, token)
		require.Equal(t, token.Type(), utils.AccessToken)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		token, err := user.GenerateRefreshToken(cfg)
		require.Nil(t, err)
		require.NotEmpty(t, token)
		require.Equal(t, token.Type(), utils.RefreshToken)
	})

	t.Run("TokenPair", func(t *testing.T) {
		tokens, err := user.GenerateTokenPair(cfg)
		require.Nil(t, err)
		require.NotEmpty(t, tokens)
		require.Equal(t, len(tokens), 2)
		require.Equal(t, tokens[0].Type(), utils.AccessToken)
		require.Equal(t, tokens[1].Type(), utils.RefreshToken)
	})
}
