package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		filepath := "."
		viperCfg, err := LoadConfig(filepath)
		require.NotNil(t, viperCfg)
		require.Nil(t, err)
	})

	t.Run("config not exists", func(t *testing.T) {
		filepath := "./nonexist"
		viperCfg, err := LoadConfig(filepath)
		require.Nil(t, viperCfg)
		require.NotNil(t, err)
	})
}

func TestParseAllowedOrigins(t *testing.T) {
	corsConfig := CorsConfig{
		AllowedOrigins: "http://example.com,https://example.org",
		Credentials:    false,
	}
	t.Run("success-origins", func(t *testing.T) {
		origins, err := corsConfig.ParseOrigins()
		require.Nil(t, err)
		require.NotNil(t, origins)
		require.Len(t, origins, 2)
		require.Equal(t, []string{"http://example.com", "https://example.org"}, origins)
	})

	t.Run("invalid-origins", func(t *testing.T) {
		cc := corsConfig
		cc.AllowedOrigins = "http//example.com,https:example.org"
		origins, err := cc.ParseOrigins()
		require.Nil(t, origins)
		require.NotNil(t, err)
	})
}
