package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRegisterRequest(t *testing.T) {
	request := RegisterRequest{
		Name:     "doe",
		Email:    "doe@mail.com",
		Password: "doe123",
	}

	t.Run("success", func(t *testing.T) {
		err := request.ValidateRequest()
		require.Nil(t, err)
	})

	t.Run("missing-name", func(t *testing.T) {
		req := request
		req.Name = ""
		err := req.ValidateRequest()
		require.NotNil(t, err)
	})

	t.Run("missing-email", func(t *testing.T) {
		req := request
		req.Email = ""
		err := req.ValidateRequest()
		require.NotNil(t, err)
	})

	t.Run("missing-password", func(t *testing.T) {
		req := request
		req.Password = ""
		err := req.ValidateRequest()
		require.NotNil(t, err)
	})
}

func TestValidateLoginRequest(t *testing.T) {
	request := LoginRequest{
		Email:    "doe@mail.com",
		Password: "doe123",
	}

	t.Run("success", func(t *testing.T) {
		err := request.ValidateRequest()
		require.Nil(t, err)
	})

	t.Run("missing-email", func(t *testing.T) {
		req := request
		req.Email = ""
		err := req.ValidateRequest()
		require.NotNil(t, err)
	})

	t.Run("missing-password", func(t *testing.T) {
		req := request
		req.Password = ""
		err := req.ValidateRequest()
		require.NotNil(t, err)
	})
}
