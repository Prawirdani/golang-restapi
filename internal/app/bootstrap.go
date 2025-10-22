package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/infra/repository/postgres"
	"github.com/prawirdani/golang-restapi/internal/service"
	"github.com/prawirdani/golang-restapi/internal/transport/http"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
)

// Init & Injects all dependencies.
func (s *Server) bootstrap() {
	// Postgres Repo Factory
	repoFact := postgres.NewRepositoryFactory(s.pg, s.logger)

	// Transactor factory
	transactor := postgres.NewTransactor(s.pg)

	// Setup Services
	authService := service.NewAuthService(
		s.cfg,
		s.logger,
		transactor,
		repoFact.User(),
		repoFact.Auth(),
	)

	// Setup Handlers
	authHandler := handler.NewAuthHandler(s.cfg, authService)

	s.router.Route("/api", func(r chi.Router) {
		http.RegisterAuthRoutes(r, authHandler, s.middlewares)
	})
}
