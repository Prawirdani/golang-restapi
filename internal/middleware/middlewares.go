package middleware

import "github.com/prawirdani/golang-restapi/config"

type MiddlewareManager struct {
	cfg *config.Config
}

func NewMiddlewareManager(cfg *config.Config) MiddlewareManager {
	return MiddlewareManager{
		cfg: cfg,
	}
}
