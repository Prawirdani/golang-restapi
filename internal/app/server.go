package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	stderrs "errors"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/app/middleware"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/logging"
	"github.com/prawirdani/golang-restapi/pkg/metrics"
	"github.com/prawirdani/golang-restapi/pkg/response"
)

type Server struct {
	cfg         *config.Config
	logger      logging.Logger
	router      *chi.Mux
	pg          *pgxpool.Pool
	middlewares *middleware.Collection
}

// Server Initialization function, also bootstraping dependency
func InitServer(cfg *config.Config, logger logging.Logger, pgPool *pgxpool.Pool) (*Server, error) {
	if cfg == nil {
		return nil, stderrs.New("Config is required")
	}

	if pgPool == nil {
		return nil, stderrs.New("Postgres connection pool is required")
	}

	router := chi.NewRouter()

	mws := middleware.NewCollection(cfg, logger)

	router.Use(mws.PanicRecoverer)
	router.Use(mws.Gzip)
	router.Use(mws.Cors)
	router.Use(mws.ReqLogger)

	if cfg.IsProduction() {
		router.Use(mws.RateLimit(50, 1*time.Minute))
	}

	// Not Found Handler
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.HandleError(w, errors.NotFound("The requested resource could not be found"))
	})

	// Request Method Not Allowed Handler
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.HandleError(
			w,
			errors.MethodNotAllowed("The method is not allowed for the requested URL"),
		)
	})

	if cfg.Metrics.Enable {
		m := metrics.Init()
		m.SetAppInfo(cfg.App.Version, string(cfg.App.Environment))
		router.Use(m.Instrument)

		// Metrics Server
		port := cfg.Metrics.PrometheusPort
		go func() {
			logger.Info(
				logging.Startup,
				"Server.Init",
				fmt.Sprintf("Metrics serves on localhost:%v", port),
			)
			if err := m.RunServer(port); err != nil {
				logger.Fatal(logging.Startup, "Server.Init.Metrics", err.Error())
			}
		}()
	}

	svr := &Server{
		router:      router,
		cfg:         cfg,
		pg:          pgPool,
		middlewares: mws,
		logger:      logger,
	}

	svr.bootstrap()

	return svr, nil
}

func (s *Server) Start() {
	fmt.Println("ENV\t:", s.cfg.App.Environment)
	fmt.Println("Metrics\t:", s.cfg.Metrics.Enable)

	svr := http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.App.Port),
		Handler: s.router,
	}

	go func() {
		s.logger.Info(
			logging.Startup,
			"app.Server",
			fmt.Sprintf("App serves on %v", s.cfg.App.Port),
		)
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(logging.Startup, "app.Server", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.logger.Info(logging.Shutdown, "app.Server.Shutdown", "Shutdown signal received")

	ctx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	if err := svr.Shutdown(ctx); err != nil {
		s.logger.Fatal(logging.Shutdown, "app.Server.Shutdown", err.Error())
	}

	s.logger.Info(logging.Shutdown, "app.Server.Shutdown", "Server gracefully stopped")
}
