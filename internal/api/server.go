package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/metrics"
)

type Server struct {
	router  *chi.Mux
	pg      *pgxpool.Pool
	metrics *metrics.Metrics
	cfg     *config.Config
}

// Server Initialization function, also bootstraping dependency
func InitServer(cfg *config.Config, pgPool *pgxpool.Pool) (*Server, error) {
	router := chi.NewRouter()

	m := metrics.Init()
	m.SetAppInfo(cfg.App.Version, string(cfg.App.Environment))

	router.Use(panicRecoverer)
	router.Use(reqLogger)
	router.Use(gzip)
	router.Use(cors(cfg))

	if cfg.IsProduction() {
		router.Use(rateLimit)
	}

	// Not Found Handler
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httputil.HandleError(w, httputil.ErrNotFound("The requested resource could not be found"))
	})

	// Request Method Not Allowed Handler
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httputil.HandleError(w, httputil.ErrMethodNotAllowed("The method is not allowed for the requested URL"))
	})

	svr := &Server{
		router:  router,
		cfg:     cfg,
		pg:      pgPool,
		metrics: m,
	}

	svr.bootstrap()

	return svr, nil
}

func (s *Server) Start() {
	svr := http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.App.Port),
		Handler: s.router,
	}

	// Application Server
	go func() {
		log.Printf("Listening on localhost%s", svr.Addr)
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed, cause: %s", err.Error())
		}
	}()

	// Metrics Server
	go s.metrics.RunServer(s.cfg.App.Port + 1)

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