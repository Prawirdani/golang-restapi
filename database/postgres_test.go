package database

import (
	"testing"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/stretchr/testify/require"
)

var c *config.Config

func init() {
	path := "../config"
	cfg, err := config.LoadConfig(path)
	if err != nil {
		panic(err)
	}
	c = cfg
}

func TestPostgresConnection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pgpool, err := NewPGConnection(c)
		require.Nil(t, err)
		require.NotNil(t, pgpool)
	})
}
