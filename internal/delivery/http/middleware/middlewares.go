package middleware

import "github.com/prawirdani/golang-restapi/config"

type MiddlewareManager struct {
	tokenCfg config.TokenConfig
}

func NewMiddlewareManager(cfg *config.Config) MiddlewareManager {
	return MiddlewareManager{
		tokenCfg: cfg.Token,
	}
}
