package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

func InitMainRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(RequestLogger)
	// Gzip Compressor
	r.Use(middleware.Compress(6))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Cors.OriginsToArray(),
		AllowedMethods:   cfg.Cors.MethodsToArray(),
		AllowCredentials: cfg.Cors.Credentials,
		Debug:            cfg.App.Environment == "dev",
	}))

	// Not Found Handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httputil.HandleError(w, httputil.ErrNotFound("The requested resource could not be found"))
	})
	// Request Method Not Allowed Handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httputil.HandleError(w, httputil.ErrMethodNotAllowed("The method is not allowed for the requested URL"))
	})

	return r
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		start := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		logAttributes := &RequestLogAttributes{
			Method:     r.Method,
			Uri:        r.RequestURI,
			ClientIP:   r.RemoteAddr,
			RequestID:  requestID,
			StatusCode: rec.Status,
			StatusText: http.StatusText(rec.Status),
			TimeTaken:  duration,
		}

		HttpRequestLogger(*logAttributes)
	})
}

type RequestLogAttributes struct {
	Method     string
	Uri        string
	ClientIP   string
	RequestID  string
	StatusCode int
	StatusText string
	TimeTaken  time.Duration
}

type ResponseRecorder struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func (rr *ResponseRecorder) WriteHeader(code int) {
	rr.Status = code
	rr.ResponseWriter.WriteHeader(code)
}

func HttpRequestLogger(rl RequestLogAttributes) {
	slog.Info(
		"HTTP Request Log",
		slog.String("method", rl.Method),
		slog.String("url", rl.Uri),
		slog.String("from", rl.ClientIP),
		slog.String("request-id", rl.RequestID),
		slog.Int("status_code", rl.StatusCode),
		slog.String("status_text", rl.StatusText),
		slog.Float64("time_taken(ms)", float64(rl.TimeTaken.Microseconds())/float64(1000)),
	)
}
