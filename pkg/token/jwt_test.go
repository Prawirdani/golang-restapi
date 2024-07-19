package token

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
		token, err := GenerateJWT(cfg, accessTokenPayload, Access)
		require.Nil(t, err)
		require.NotEmpty(t, token)

		require.Equal(t, token.Type(), Access)
		require.Equal(t, token.Payload()["id"], accessTokenPayload["id"])
		require.Equal(t, token.Payload()["name"], accessTokenPayload["name"])
	})

	t.Run("Generate-RefreshToken", func(t *testing.T) {
		rfToken, err := GenerateJWT(cfg, refreshTokenPayload, Refresh)
		require.Nil(t, err)
		require.NotEmpty(t, rfToken)
		require.Equal(t, rfToken.Type(), Refresh)
		require.Equal(t, rfToken.Payload()["id"], refreshTokenPayload["id"])
	})
}

func TestParseToken(t *testing.T) {
	t.Run("Parse-AccessToken", func(t *testing.T) {
		token, err := GenerateJWT(cfg, accessTokenPayload, Access)
		require.Nil(t, err)
		require.NotEmpty(t, token)

		parsedToken, err := parseJWT(token.String(), cfg.Token.SecretKey)
		require.Nil(t, err)
		require.NotEmpty(t, parsedToken)
		require.Equal(t, parsedToken["user"], accessTokenPayload)
	})

	t.Run("Parse-RefreshToken", func(t *testing.T) {
		rfToken, err := GenerateJWT(cfg, refreshTokenPayload, Refresh)
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

		acToken, err := GenerateJWT(modifiedCfg, accessTokenPayload, Access)
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

		rfToken, err := GenerateJWT(modifiedCfg, accessTokenPayload, Access)
		require.Nil(t, err)
		require.NotEmpty(t, rfToken)

		parsedToken, err := parseJWT(rfToken.String(), cfg.Token.SecretKey)
		require.NotNil(t, err)
		require.Nil(t, parsedToken)
		require.Equal(t, err, ErrorTokenInvalid)
	})
}