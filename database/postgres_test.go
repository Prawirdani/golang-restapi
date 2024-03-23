package database

import (
	"testing"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/stretchr/testify/require"
)

var c *config.Config

func init() {
	path := "../config"
	viper := config.LoadConfig(path)
	c = config.ParseConfig(viper)
}

func TestPostgresConnection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pgpool, err := NewPGConnection(c.DB)
		require.Nil(t, err)
		require.NotNil(t, pgpool)
	})
}
