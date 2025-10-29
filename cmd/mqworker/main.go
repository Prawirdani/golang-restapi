package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	"github.com/prawirdani/golang-restapi/internal/infra/mq/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/infra/mq/worker"
	"github.com/prawirdani/golang-restapi/internal/mail"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		stdlog.Fatal("Failed to load config", err)
	}
	log.SetLogger(log.NewZerologAdapter(cfg))

	mailer := mail.NewMailer(cfg)

	emailWorker, err := worker.NewEmailWorker(mailer)
	if err != nil {
		log.Error("Failed to init email worker", "err", err.Error())
		os.Exit(1)
	}

	rmqConsumer, err := rabbitmq.NewConsumer(cfg.RabbitMqURL)
	if err != nil {
		log.Error("Failed to init rabbitmq consumer", "err", err.Error())
		os.Exit(1)
	}
	defer rmqConsumer.Close()

	rmqConsumer.RegisterHandler(mq.EmailResetPasswordJobKey, emailWorker.HandlePasswordReset)

	// Start consuming
	ctx := context.Background()
	log.Info("RabbitMQ consumer started")
	rmqConsumer.Start(ctx)
}
