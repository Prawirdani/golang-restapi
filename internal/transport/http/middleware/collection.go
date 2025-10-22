package middleware

import (
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type Collection struct {
	cfg    *config.Config
	logger logging.Logger
}

func NewCollection(cfg *config.Config, logger logging.Logger) *Collection {
	mw := Collection{
		cfg:    cfg,
		logger: logger,
	}

	return &mw
}
