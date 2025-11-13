package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/infra/messaging/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/infra/repository/postgres"
	"github.com/prawirdani/golang-restapi/internal/infra/storage/r2"
	"github.com/prawirdani/golang-restapi/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Services struct {
	UserService *service.UserService
	AuthService *service.AuthService
}

// Container holds all application dependencies
type Container struct {
	Config   *config.Config
	Services *Services
	pgpool   *pgxpool.Pool
}

// NewContainer initializes all dependencies
func NewContainer(
	cfg *config.Config,
	pgpool *pgxpool.Pool,
	rmqconn *amqp.Connection,
) (*Container, error) {
	// Postgres Repo Factory
	repoFactory := postgres.NewRepositoryFactory(pgpool)
	transactor := postgres.NewTransactor(pgpool)

	r2PublicStorage, err := r2.New(r2.Config{
		BucketURL:       cfg.R2.PublicBucketURL,
		BucketName:      cfg.R2.PublicBucket,
		AccountID:       cfg.R2.AccountID,
		AccessKeyID:     cfg.R2.AccessKeyID,
		AccessKeySecret: cfg.R2.AccessKeySecret,
	})
	if err != nil {
		return nil, err
	}

	// Setup Services
	userService := service.NewUserService(transactor, repoFactory.User(), r2PublicStorage)

	authMessagePublisher := rabbitmq.NewAuthMessagePublisher(rmqconn)
	authService := service.NewAuthService(
		cfg.Auth,
		transactor,
		repoFactory.User(),
		repoFactory.Auth(),
		authMessagePublisher,
	)

	c := &Container{
		Config: cfg,
		Services: &Services{
			UserService: userService,
			AuthService: authService,
		},
		pgpool: pgpool,
	}

	return c, nil
}
