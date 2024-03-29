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
