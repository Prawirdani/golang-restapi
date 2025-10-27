package middleware

import (
	"github.com/prawirdani/golang-restapi/config"
)

type Collection struct {
	cfg *config.Config
}

func Setup(cfg *config.Config) *Collection {
	mw := Collection{cfg: cfg}

	return &mw
}
