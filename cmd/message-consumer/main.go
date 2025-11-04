package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	stdlog "log"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/infra/messaging/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/transport/amqp/consumer"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/mailer"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		stdlog.Fatal("Failed to load config", err)
	}
	log.SetLogger(log.NewZerologAdapter(cfg))

	conn, err := rabbitmq.Dial(cfg.RabbitMqURL)
	if err != nil {
		log.Error("Failed to dial rabbitmq connection", "error", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Channel for setup topology
	ch, err := conn.Channel()
	if err != nil {
		log.Error("Failed to create rmq channel", "error", err.Error())
		os.Exit(1)
	}

	if err := rabbitmq.SetupTopology(ch); err != nil {
		log.Error("Failed to setup rmq topology", "error", err.Error())
		os.Exit(1)
	}
	ch.Close()

	// mailQueueHandlers := queueHandler.NewAuthMessageConsumer(m)
	consumerRegistry := consumer.NewRegistry(conn)

	m := mailer.New(cfg)
	authConsumers := consumer.NewAuthMessageConsumer(m)

	consumerRegistry.RegisterConsumers(map[string]consumer.HandlerFunc{
		rabbitmq.AuthEmailResetPasswordQueue.Name: authConsumers.EmailResetPasswordHandler,
	})

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()
	if err := consumerRegistry.Start(ctx); err != nil && err != context.Canceled {
		log.Error("Message consumer exited", "error", err.Error())
	}

	log.Info("Message consumer exited gracefully")
}
