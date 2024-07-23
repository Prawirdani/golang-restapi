package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/handler"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/internal/service"
)

// Init & Injects all dependencies.
func (s *Server) bootstrap() {
	// Setup Repos
	userRepository := repository.NewUserRepository(s.pg, s.logger)

	// Setup Services
	authUC := service.NewAuthService(s.cfg, s.logger, userRepository)

	// Setup Handlers
	authHandler := handler.NewAuthHandler(s.cfg, authUC)

	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.metrics.Instrument)
		// route.RegisterAuthRoutes(r, authHandler, mws)
		registerAuthRoutes(r, authHandler, s.middlewares)
	})
}
