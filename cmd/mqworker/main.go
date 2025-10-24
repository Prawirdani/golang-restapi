package main

import (
	"context"
	"log"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	"github.com/prawirdani/golang-restapi/internal/infra/mq/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/infra/mq/worker"
	"github.com/prawirdani/golang-restapi/internal/mail"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(cfg)
	mailer := mail.NewMailer(cfg, logger)

	emailWorker, err := worker.NewEmailWorker(mailer, logger)
	if err != nil {
		logger.Fatal(logging.MQWorker, "main.NewEmailWorker", err.Error())
	}

	rmqConsumer, err := rabbitmq.NewConsumer(cfg.RabbitMqURL)
	if err != nil {
		logger.Fatal(logging.MQWorker, "main.NewConsumer", err.Error())
	}
	defer rmqConsumer.Close()

	rmqConsumer.RegisterHandler(mq.EmailResetPasswordJobKey, emailWorker.HandlePasswordReset)

	// Start consuming
	ctx := context.Background()
	logger.Info(logging.MQWorker, "main.Start", "MQ Consumer Started")
	if err := rmqConsumer.Start(ctx); err != nil {
		logger.Fatal(logging.MQWorker, "main.Start", err.Error())
	}
}
