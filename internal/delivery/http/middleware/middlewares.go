package middleware

import "github.com/prawirdani/golang-restapi/config"

type Collection struct {
	tokenCfg config.TokenConfig
}

func New(cfg *config.Config) Collection {
	return Collection{
		tokenCfg: cfg.Token,
	}
}
