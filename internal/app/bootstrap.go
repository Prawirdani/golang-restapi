package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/infra/mq/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/infra/repository/postgres"
	"github.com/prawirdani/golang-restapi/internal/infra/storage/r2"
	"github.com/prawirdani/golang-restapi/internal/service"
	"github.com/prawirdani/golang-restapi/internal/transport/http"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

// Init & Injects all dependencies.
func (s *Server) bootstrap() {
	// Postgres Repo Factory
	repoFactory := postgres.NewRepositoryFactory(s.pg, s.logger)

	// Transactor factory
	transactor := postgres.NewTransactor(s.pg)

	publicR2, err := r2.New(r2.Config{
		BucketURL:       s.cfg.R2.PublicBucketURL,
		BucketName:      s.cfg.R2.PublicBucket,
		AccountID:       s.cfg.R2.AccountID,
		AccessKeyID:     s.cfg.R2.AccessKeyID,
		AccessKeySecret: s.cfg.R2.AccessKeySecret,
	})
	if err != nil {
		s.logger.Fatal(logging.Startup, "Server.bootstrap", err.Error())
	}

	rmqProducer, err := rabbitmq.NewPublisher(s.cfg.RabbitMqURL)
	if err != nil {
		rmqProducer.Close()
		s.logger.Fatal(logging.Startup, "Server.bootstrap", err.Error())
	}

	// Setup Services
	userService := service.NewUserService(
		s.cfg,
		s.logger,
		transactor,
		repoFactory.User(),
		publicR2,
	)
	authService := service.NewAuthService(
		s.cfg,
		s.logger,
		transactor,
		repoFactory.User(),
		repoFactory.Auth(),
		userService,
		rmqProducer,
	)

	// Setup Handlers
	userHandler := handler.NewUserHandler(s.logger, userService)
	authHandler := handler.NewAuthHandler(s.cfg, authService)

	s.router.Route("/api", func(r chi.Router) {
		http.RegisterUserROutes(r, userHandler, s.middlewares)
		http.RegisterAuthRoutes(r, authHandler, s.middlewares)
	})
}
