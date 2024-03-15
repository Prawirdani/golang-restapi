package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prawirdani/golang-restapi/internal/app"
	"github.com/spf13/viper"
)

func main() {
	cfg := app.NewConfig()
	app.InitLogger(cfg)
	jsonHandler := app.NewJsonHandler()

	dbpool := app.NewPGPool(cfg)
	router := app.InitMainRouter(cfg, jsonHandler)

	bootstrap := app.Configuration{
		MainRouter: router,
		Config:     cfg,
		DBPool:     dbpool,
		JSON:       jsonHandler,
	}
	app.Bootstrap(&bootstrap)

	server := NewServer(cfg, router)
	server.Start()
}

type Server struct {
	*http.Server
}

func NewServer(config *viper.Viper, multiplexer http.Handler) *Server {
	svr := http.Server{
		Addr:    fmt.Sprintf(":%v", config.GetInt("app.port")),
		Handler: multiplexer,
	}

	return &Server{&svr}
}

func (s *Server) Start() {
	go func() {
		log.Printf("Listening on 0.0.0.0%s", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed, cause: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed, cause: %s", err.Error())
	}

	log.Println("Server gracefully stopped")
}
