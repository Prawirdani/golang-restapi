package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/delivery/http"
	"github.com/prawirdani/golang-restapi/internal/middleware"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/internal/usecase"
)

type Configuration struct {
	MainRouter *chi.Mux
	DBPool     *pgxpool.Pool
	Config     config.Config
}

// Init & Injects all dependencies.
// This function should be called at main.go file to set up all required services and components.
func Bootstrap(c Configuration) {
	// Setup Repos
	userRepository := repository.NewUserRepository(c.DBPool, "users")

	// Setup Usecases
	authUC := usecase.NewAuthUseCase(c.Config, userRepository)

	// Setup Handlers
	authHandler := http.NewAuthHandler(c.Config, authUC)

	middlewares := middleware.NewMiddlewareManager(c.Config)

	c.MainRouter.Route("/api/v1", func(v1 chi.Router) {
		http.MapAuthRoutes(v1, authHandler, middlewares)
	})
}
