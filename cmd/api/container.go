package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/domain/auth"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/internal/infrastructure/messaging/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/infrastructure/repository/postgres"
	"github.com/prawirdani/golang-restapi/internal/infrastructure/storage/r2"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Services struct {
	UserService *user.Service
	AuthService *auth.Service
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
	userService := user.NewService(transactor, repoFactory.User(), r2PublicStorage)

	authMessagePublisher := rabbitmq.NewAuthMessagePublisher(rmqconn)
	authService := auth.NewService(
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
