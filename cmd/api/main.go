package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	stdlog "log"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/infra/messaging/rabbitmq"
	"github.com/prawirdani/golang-restapi/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		stdlog.Fatal("Failed to load config", err)
	}
	log.SetLogger(log.NewZerologAdapter(cfg))

	pgpool, err := database.NewPGConnection(cfg)
	if err != nil {
		log.Error("Failed to create postgres connection", "error", err.Error())
		os.Exit(1)
	}
	defer pgpool.Close()

	rmqconn, err := initRabbitMQ(cfg.RabbitMqURL)
	if err != nil {
		log.Error("Failed to init rabbit mq", "error", err.Error())
		os.Exit(1)
	}
	defer rmqconn.Close()

	container, err := NewContainer(cfg, pgpool, rmqconn)
	if err != nil {
		log.Error("Failed to create container", "error", err.Error())
		os.Exit(1)
	}

	server, err := NewServer(container)
	if err != nil {
		log.Error("Failed to create server", "error", err.Error())
		os.Exit(1)
	}

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

	// Start message consumers in goroutine because blocking
	go func() {
		if err := startMessageConsumers(ctx, rmqconn, cfg); err != nil && err != context.Canceled {
			log.Error("Worker exited with error", "error", err.Error())
			cancel()
		}
	}()

	// Run HTTP server
	if err := server.Start(ctx); err != nil {
		log.Error("Server exited with error", "error", err.Error())
	}

	log.Info("Application exited gracefully")
}

func initRabbitMQ(url string) (*amqp.Connection, error) {
	conn, err := rabbitmq.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	if err := rabbitmq.SetupTopologies(
		conn,
		rabbitmq.ResetPasswordEmailTopology,
	); err != nil {
		return nil, fmt.Errorf("setup topologies: %w", err)
	}

	return conn, nil
}
