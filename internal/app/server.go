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
	res "github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/metrics"

	httptransport "github.com/prawirdani/golang-restapi/internal/transport/http"
)

const MAX_BODY_SIZE = 10 << 20 // 10MB

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

	mws := middleware.Setup(container.Config)

	// Apply global middlewares
	router.Use(middleware.RequestID)
	router.Use(mws.MaxBodySizeMiddleware(MAX_BODY_SIZE))
	router.Use(mws.PanicRecovery)
	router.Use(
		mws.Gzip,
	) // TODO: Should based on config, since proxy mostly able to handle compression
	router.Use(mws.Cors)
	router.Use(mws.ReqLogger)

	if container.Config.IsProduction() {
		router.Use(mws.RateLimit(50, 1*time.Minute))
	}

	// Error handlers
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res.HandleError(w, errors.NotFound("The requested resource could not be found"))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		res.HandleError(
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

	router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		res.Send(w, r, res.WithMessage("services up and running"))
	})

	return svr, nil
}

func (s *Server) Start() {
	cfg := s.container.Config

	fmt.Println("Environment\t:", cfg.App.Environment)
	fmt.Println("Metrics\t\t:", cfg.Metrics.Enable)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.App.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Info(fmt.Sprintf("server listening on 0.0.0.0:%v", cfg.App.Port))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", "error", err.Error())
	}

	// Cleanup resources
	s.container.Cleanup()

	log.Info("server stopped gracefully")
}

func (s *Server) setupMetrics() {
	i := metrics.Init(
		s.container.Config.App.Version,
		string(s.container.Config.App.Environment),
	)
	s.router.Use(i.Instrument)

	// Metrics server
	port := s.container.Config.Metrics.PrometheusPort
	go func() {
		if err := i.RunServer(port); err != nil {
			log.Error("failed to run metrics server", "err", err.Error())
		}
		log.Info(fmt.Sprintf("metrics serves on 0.0.0.0:%v/metrics", port))
	}()
}

func (s *Server) setupHandlers() {
	svcs := s.container.Services
	// Setup Handlers
	userHandler := handler.NewUserHandler(svcs.UserService)
	authHandler := handler.NewAuthHandler(s.container.Config, svcs.AuthService)

	s.router.Route("/api", func(r chi.Router) {
		httptransport.RegisterUserRoutes(r, userHandler, s.middlewares)
		httptransport.RegisterAuthRoutes(r, authHandler, s.middlewares)
	})
}
