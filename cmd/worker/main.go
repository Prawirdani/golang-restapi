package main

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/worker"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/mailer"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		stdlog.Fatal("Failed to load config", err)
	}
	log.SetLogger(log.NewZerologAdapter(cfg.IsProduction()))

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0, // use default DB
	})
	defer rdb.Close()

	mailer := mailer.New(cfg.SMTP)
	emailEventConsumer := worker.NewEmailEventConsumer(mailer).Handler(rdb)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	// Run the consumer and surface its exit. Start blocks until ctx is cancelled
	// (returning context.Canceled after draining in-flight handlers) or it hits a
	// fatal error.
	done := make(chan error, 1)
	go func() {
		done <- emailEventConsumer.Start(ctx)
	}()

	select {
	case <-ctx.Done():
		// Shutdown requested: wait for Start to return so in-flight emails drain
		// and deferred cleanup (rdb.Close) runs.
		if err := <-done; err != nil && !errors.Is(err, context.Canceled) {
			log.Error("Consumer stopped with error during shutdown", err)
		}
		log.Info("Worker exited gracefully")
	case err := <-done:
		// Consumer returned on its own before a shutdown signal.
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Error("Consumer stopped unexpectedly", err)
			cancel()
			os.Exit(1)
		}
		log.Info("Worker exited gracefully")
	}
}
