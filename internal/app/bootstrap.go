package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/delivery/http"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/middleware"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/internal/usecase"
)

type Configuration struct {
	MainRouter *chi.Mux
	DBPool     *pgxpool.Pool
}

// Init & Injects all dependencies.
// This function should be called at main.go file to set up all required services and components.
func Bootstrap(c *Configuration) {

	// Setup Repos
	userRepository := repository.NewUserRepository(c.DBPool, "users")

	// Setup Usecases
	authUC := usecase.NewAuthUseCase(userRepository)

	middlewares := middleware.New()
	// Setup Handlers
	authHandler := http.NewAuthHandler(middlewares, authUC)

	routes := http.SetupAPIRoutes(c.MainRouter)
	routes.RegisterHandlers(authHandler)
	routes.Init()
}
