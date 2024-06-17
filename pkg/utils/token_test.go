package utils

import (
	"log"
	"testing"

	"github.com/prawirdani/golang-restapi/config"
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

var (
	accessTokenPayload = map[string]interface{}{
		"id":   "user-id-1",
		"name": "doe",
	}
	refreshTokenPayload = map[string]interface{}{
		"id": "user-id-1",
	}
)

func TestGenerateJWT(t *testing.T) {
	t.Run("Generate-AccessToken", func(t *testing.T) {
		acToken, err := GenerateJWT(cfg, accessTokenPayload, AccessToken)
		require.Nil(t, err)
		require.NotEmpty(t, acToken)

		require.Equal(t, acToken.Claims.TokenType, AccessToken)
		require.Equal(t, acToken.Claims.User["id"], accessTokenPayload["id"])
		require.Equal(t, acToken.Claims.User["name"], accessTokenPayload["name"])
	})

	t.Run("Generate-RefreshToken", func(t *testing.T) {
		rfToken, err := GenerateJWT(cfg, refreshTokenPayload, RefreshToken)
		require.Nil(t, err)
		require.NotEmpty(t, rfToken)
		require.Equal(t, rfToken.Claims.TokenType, RefreshToken)
		require.Equal(t, rfToken.Claims.User["id"], refreshTokenPayload["id"])
	})
}

func TestParseToken(t *testing.T) {
	t.Run("Parse-AccessToken", func(t *testing.T) {
		acToken, err := GenerateJWT(cfg, accessTokenPayload, AccessToken)
		require.Nil(t, err)
		require.NotEmpty(t, acToken)

		parsedToken, err := parseJWT(acToken.String(), cfg.Token.SecretKey)
		require.Nil(t, err)
		require.NotEmpty(t, parsedToken)
		require.Equal(t, parsedToken["user"], accessTokenPayload)
	})

	t.Run("Parse-RefreshToken", func(t *testing.T) {
		rfToken, err := GenerateJWT(cfg, refreshTokenPayload, RefreshToken)
		require.Nil(t, err)
		require.NotEmpty(t, rfToken)

		parsedToken, err := parseJWT(rfToken.String(), cfg.Token.SecretKey)
		require.Nil(t, err)
		require.NotEmpty(t, parsedToken)
		require.Equal(t, parsedToken["user"], refreshTokenPayload)
	})

	t.Run("Expired-AccessToken", func(t *testing.T) {
		modifiedCfg := cfg
		modifiedCfg.Token.AccessTokenExpiry = -5

		acToken, err := GenerateJWT(modifiedCfg, accessTokenPayload, AccessToken)
		require.Nil(t, err)
		require.NotEmpty(t, acToken)

		parsedToken, err := parseJWT(acToken.String(), cfg.Token.SecretKey)
		require.NotNil(t, err)
		require.Nil(t, parsedToken)
		require.Equal(t, err, ErrorTokenInvalid)
	})

	t.Run("Expired-RefreshToken", func(t *testing.T) {
		modifiedCfg := cfg
		modifiedCfg.Token.RefreshTokenExpiry = -5

		rfToken, err := GenerateJWT(modifiedCfg, accessTokenPayload, AccessToken)
		require.Nil(t, err)
		require.NotEmpty(t, rfToken)

		parsedToken, err := parseJWT(rfToken.String(), cfg.Token.SecretKey)
		require.NotNil(t, err)
		require.Nil(t, parsedToken)
		require.Equal(t, err, ErrorTokenInvalid)
	})
}
