package main

import (
	"context"
	"fmt"
	"net/http"
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
	metrics     *metrics.Metrics
	middlewares *middleware.Collection
}

// NewServer acts as a constructor, initializing the server and its dependencies.
// All router setup is deferred to the setupRouter method.
func NewServer(container *Container) (*Server, error) {
	if container == nil {
		return nil, fmt.Errorf("container is required")
	}

	svr := &Server{
		container: container,
		router:    chi.NewRouter(),
		metrics: metrics.Init(
			container.Config.App.Version,
			string(container.Config.App.Environment),
			container.Config.Metrics.PrometheusPort,
		),
		middlewares: middleware.Setup(container.Config),
	}

	return svr, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Configure all routes, middlewares, and handlers
	s.setupRouter()

	cfg := s.container.Config

	// Metrics server
	var metricServer *http.Server
	if cfg.IsProduction() {
		metricServer = &http.Server{
			Addr:    fmt.Sprintf(":%v", cfg.Metrics.PrometheusPort),
			Handler: s.metrics.ExporterHandler(),
		}

		// Start metrics server
		go func() {
			log.Info(
				fmt.Sprintf("Metrics serving on 0.0.0.0:%v/metrics", cfg.Metrics.PrometheusPort),
			)
			if err := metricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error("Metrics server stopped unexpectedly", err)
			}
		}()
	}

	// API server
	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.App.Port),
		Handler:      s.router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start API server
	go func() {
		log.Info(fmt.Sprintf("API server listening on 0.0.0.0:%v", cfg.App.Port))
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("API server stopped unexpectedly", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Failed to shutdown API server", err)
	}

	if metricServer != nil {
		if err := metricServer.Shutdown(shutdownCtx); err != nil {
			log.Error("Failed to shutdown Metrics server", err)
		}
	}
	return nil
}

// setupRouter configures all middlewares, error handlers, and API routes.
func (s *Server) setupRouter() {
	mws := s.middlewares

	if s.container.Config.IsProduction() {
		s.router.Use(middleware.RequestID)
		s.router.Use(mws.RateLimit(50, 1*time.Minute))
		s.router.Use(s.metrics.InstrumentHandler) // Instrument the main router
	} else {
		s.router.Use(mws.ReqLogger)
	}

	// Apply common middlewares
	s.router.Use(mws.MaxBodySizeMiddleware(MAX_BODY_SIZE))
	s.router.Use(mws.PanicRecovery)
	s.router.Use(mws.Gzip) // TODO: Should be based on config
	s.router.Use(mws.Cors)

	// Custom 404 and 405 handlers
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res.HandleError(w, errors.NotFound("The requested resource could not be found"))
	})
	s.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		res.HandleError(
			w,
			errors.MethodNotAllowed("The method is not allowed for the requested URL"),
		)
	})

	// Setup API routes
	s.setupHandlers()

	// Health check route
	s.router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		res.Send(w, r, res.WithMessage("services up and running"))
	})
}

// setupHandlers initializes and registers all API handlers.
func (s *Server) setupHandlers() {
	svcs := s.container.Services

	// Initialize Handlers
	userHandler := handler.NewUserHandler(svcs.UserService)
	authHandler := handler.NewAuthHandler(s.container.Config, svcs.AuthService)

	// Register API routes
	s.router.Route("/api", func(r chi.Router) {
		httptransport.RegisterUserRoutes(r, userHandler, s.middlewares)
		httptransport.RegisterAuthRoutes(r, authHandler, s.middlewares)
	})
}
