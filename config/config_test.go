package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		viperCfg, err := LoadConfig()
		require.NotNil(t, viperCfg)
		require.Nil(t, err)
	})
}
