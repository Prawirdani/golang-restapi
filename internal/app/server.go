package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type Server struct {
	pg     *pgxpool.Pool
	router *chi.Mux
	cfg    config.Config
}

// Server Initialization function, also bootstraping dependency
func InitServer(cfg config.Config, pgPool *pgxpool.Pool) (*Server, error) {
	router := chi.NewRouter()

	logger := httplog.NewLogger("request-logger", httplog.Options{
		LogLevel:         slog.LevelDebug,
		JSON:             cfg.IsProduction(),
		Concise:          cfg.IsProduction(),
		RequestHeaders:   true,
		ResponseHeaders:  true,
		MessageFieldName: "message",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"version": cfg.App.Version,
			"env":     string(cfg.App.Environment),
		},
		QuietDownRoutes: []string{"/"},
		QuietDownPeriod: 10 * time.Second,
	})

	router.Use(httplog.RequestLogger(logger))

	// Gzip Compressor
	router.Use(middleware.Compress(6))

	origins, err := cfg.Cors.ParseOrigins()
	if err != nil {
		return nil, err
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "HEAD", "DELETE", "OPTIONS"},
		AllowCredentials: cfg.Cors.Credentials,
		Debug:            cfg.IsProduction(),
	}))

	// Not Found Handler
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httputil.HandleError(w, httputil.ErrNotFound("The requested resource could not be found"))
	})
	// Request Method Not Allowed Handler
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httputil.HandleError(w, httputil.ErrMethodNotAllowed("The method is not allowed for the requested URL"))
	})

	svr := &Server{
		router: router,
		cfg:    cfg,
		pg:     pgPool,
	}

	svr.bootstrap()

	return svr, nil
}

func (s *Server) Start() {

	svr := http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.App.Port),
		Handler: s.router,
	}

	go func() {
		log.Printf("Listening on 0.0.0.0%s", svr.Addr)
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed, cause: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed, cause: %s", err.Error())
	}

	log.Println("Server gracefully stopped")
}
