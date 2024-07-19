package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/delivery/http"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/middleware"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/internal/service"
)

// Init & Injects all dependencies.
func (s *Server) bootstrap() {
	// Setup Repos
	userRepository := repository.NewUserRepository(s.pg)

	// Setup Services
	authUC := service.NewAuthService(s.cfg, userRepository)

	// Setup Handlers
	authHandler := http.NewAuthHandler(s.cfg, authUC)

	mws := middleware.NewCollection(s.cfg)
	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.metrics.Instrument)
		http.RegisterAuthRoutes(r, authHandler, mws)

	})
}
