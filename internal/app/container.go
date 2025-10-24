package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	"github.com/prawirdani/golang-restapi/internal/infra/mq/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/infra/repository/postgres"
	"github.com/prawirdani/golang-restapi/internal/infra/storage/r2"
	"github.com/prawirdani/golang-restapi/internal/service"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type Services struct {
	UserService *service.UserService
	AuthService *service.AuthService
}

// Container holds all application dependencies
type Container struct {
	Config     *config.Config
	Logger     logging.Logger
	Services   *Services
	pgpool     *pgxpool.Pool
	mqproducer mq.MessageProducer
}

// NewContainer initializes all dependencies
func NewContainer(cfg *config.Config) *Container {
	logger := logging.NewLogger(cfg)

	pgpool, err := database.NewPGConnection(cfg)
	if err != nil {
		logger.Fatal(logging.Postgres, "main.NewPGConnection", err.Error())
	}

	// Postgres Repo Factory
	repoFactory := postgres.NewRepositoryFactory(pgpool, logger)
	transactor := postgres.NewTransactor(pgpool)

	rmqproducer, err := rabbitmq.NewPublisher(cfg.RabbitMqURL)
	if err != nil {
		logger.Fatal(logging.Startup, "Server.bootstrap", err.Error())
	}

	r2PublicStorage, err := r2.New(r2.Config{
		BucketURL:       cfg.R2.PublicBucketURL,
		BucketName:      cfg.R2.PublicBucket,
		AccountID:       cfg.R2.AccountID,
		AccessKeyID:     cfg.R2.AccessKeyID,
		AccessKeySecret: cfg.R2.AccessKeySecret,
	})
	if err != nil {
		logger.Fatal(logging.Startup, "Server.bootstrap", err.Error())
	}

	// Setup Services
	userService := service.NewUserService(
		cfg,
		logger,
		transactor,
		repoFactory.User(),
		r2PublicStorage,
	)
	authService := service.NewAuthService(
		cfg,
		logger,
		transactor,
		repoFactory.User(),
		repoFactory.Auth(),
		userService,
		rmqproducer,
	)

	c := &Container{
		Config: cfg,
		Logger: logger,
		Services: &Services{
			UserService: userService,
			AuthService: authService,
		},
		pgpool:     pgpool,
		mqproducer: rmqproducer,
	}

	return c
}

func (c *Container) Cleanup() error {
	if c.pgpool != nil {
		c.pgpool.Close()
	}

	if c.mqproducer != nil {
		c.mqproducer.Close()
	}

	return nil
}
