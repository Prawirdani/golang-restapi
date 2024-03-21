package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/delivery/http"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/middleware"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/internal/usecase"
	"github.com/prawirdani/golang-restapi/pkg/utils"
	"github.com/spf13/viper"
)

type Configuration struct {
	MainRouter *chi.Mux
	DBPool     *pgxpool.Pool
	Config     *viper.Viper
}

// Init & Injects all dependencies.
// This function should be called at main.go file to set up all required services and components.
func Bootstrap(c *Configuration) {
	jwtProvider := utils.NewJWTProvider(c.Config)

	// Setup Repos
	userRepository := repository.NewUserRepository("users")

	// Setup Usecases
	authUC := usecase.NewAuthUseCase(c.DBPool, userRepository, jwtProvider)

	middlewares := middleware.New(jwtProvider)
	// Setup Handlers
	authHandler := http.NewAuthHandler(middlewares, authUC)

	routes := http.SetupAPIRoutes(c.MainRouter)
	routes.RegisterHandlers(authHandler)
	routes.Init()
}
