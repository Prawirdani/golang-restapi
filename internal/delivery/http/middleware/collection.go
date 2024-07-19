package middleware

import (
	"github.com/prawirdani/golang-restapi/config"
)

type Collection struct {
	cfg *config.Config
}

func NewCollection(cfg *config.Config) *Collection {
	mw := Collection{
		cfg: cfg,
	}

	return &mw
}
