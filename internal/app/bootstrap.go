package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/delivery/http"
	"github.com/prawirdani/golang-restapi/internal/middleware"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/internal/usecase"
)

// Init & Injects all dependencies.
func (s Server) bootstrap() {
	// Setup Repos
	userRepository := repository.NewUserRepository(s.pg, "users")

	// Setup Usecases
	authUC := usecase.NewAuthUseCase(s.cfg, userRepository)

	// Setup Handlers
	authHandler := http.NewAuthHandler(s.cfg, authUC)

	middlewares := middleware.NewMiddlewareManager(s.cfg)

	s.router.Route("/api/v1", func(v1 chi.Router) {
		http.MapAuthRoutes(v1, authHandler, middlewares)
	})
}
