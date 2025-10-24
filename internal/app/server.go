package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
	"github.com/prawirdani/golang-restapi/internal/transport/http/middleware"
	"github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/logging"
	"github.com/prawirdani/golang-restapi/pkg/metrics"

	httptransport "github.com/prawirdani/golang-restapi/internal/transport/http"
)

type Server struct {
	container   *Container
	router      *chi.Mux
	middlewares *middleware.Collection
}

func NewServer(container *Container) (*Server, error) {
	if container == nil {
		return nil, fmt.Errorf("container is required")
	}

	router := chi.NewRouter()
	mws := middleware.NewCollection(container.Config, container.Logger)

	// Apply global middlewares
	router.Use(mws.PanicRecoverer)
	router.Use(mws.Gzip)
	router.Use(mws.Cors)
	router.Use(mws.ReqLogger)

	if container.Config.IsProduction() {
		router.Use(mws.RateLimit(50, 1*time.Minute))
	}

	// Error handlers
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.HandleError(w, errors.NotFound("The requested resource could not be found"))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.HandleError(
			w,
			errors.MethodNotAllowed("The method is not allowed for the requested URL"),
		)
	})

	svr := &Server{
		container:   container,
		router:      router,
		middlewares: mws,
	}

	// Setup metrics if enabled
	if container.Config.Metrics.Enable {
		svr.setupMetrics()
	}

	// Setup routes
	svr.setupHandlers()

	return svr, nil
}

func (s *Server) Start() {
	cfg := s.container.Config
	logger := s.container.Logger

	fmt.Println("ENV\t:", cfg.App.Environment)
	fmt.Println("Metrics\t:", cfg.Metrics.Enable)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.App.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		logger.Info(
			logging.Startup,
			"Server.Start",
			fmt.Sprintf("App serves on %v", cfg.App.Port),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(logging.Startup, "Server.Start", err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info(logging.Shutdown, "Server.Shutdown", "Shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal(logging.Shutdown, "Server.Shutdown", err.Error())
	}

	// Cleanup resources
	s.container.Cleanup()

	logger.Info(logging.Shutdown, "Server.Shutdown", "Server gracefully stopped")
}

func (s *Server) setupMetrics() {
	m := metrics.Init()
	m.SetAppInfo(
		s.container.Config.App.Version,
		string(s.container.Config.App.Environment),
	)
	s.router.Use(m.Instrument)

	// Metrics server
	port := s.container.Config.Metrics.PrometheusPort
	go func() {
		s.container.Logger.Info(
			logging.Startup,
			"Server.setupMetrics",
			fmt.Sprintf("Metrics serves on localhost:%v", port),
		)
		if err := m.RunServer(port); err != nil {
			s.container.Logger.Fatal(logging.Startup, "Server.setupMetrics", err.Error())
		}
	}()
}

func (s *Server) setupHandlers() {
	svcs := s.container.Services
	// Setup Handlers
	userHandler := handler.NewUserHandler(s.container.Logger, svcs.UserService)
	authHandler := handler.NewAuthHandler(s.container.Logger, s.container.Config, svcs.AuthService)

	s.router.Route("/api", func(r chi.Router) {
		httptransport.RegisterUserRoutes(r, userHandler, s.middlewares)
		httptransport.RegisterAuthRoutes(r, authHandler, s.middlewares)
	})
}
