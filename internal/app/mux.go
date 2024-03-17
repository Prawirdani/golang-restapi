package app

import (
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/spf13/viper"
)

func InitMainRouter(config *viper.Viper) *chi.Mux {
	r := chi.NewRouter()

	r.Use(RequestLogger)
	// Gzip Compressor
	r.Use(middleware.Compress(6))

	allowedOrigins := func() []string {
		origins := strings.Split(config.GetString("cors.origins"), ",")
		// Validate Origins URL
		for _, origin := range origins {
			_, err := url.ParseRequestURI(origin)
			if err != nil {
				log.Fatal(err)
			}
		}
		return origins
	}()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   strings.Split(config.GetString("cors.methods"), ","),
		AllowCredentials: config.GetBool("cors.credentials"),
		Debug:            config.GetString("app.env") == "dev",
	}))

	// Not Found Handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httputil.SendError(w, httputil.ErrNotFound("The requested url is not found"))
	})
	// Request Method Not Allowed Handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httputil.SendError(w, httputil.ErrMethodNotAllowed("The method is not allowed for the requested URL"))
	})

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		res := httputil.Response{Message: "Service up and running"}
		httputil.SendJson(w, http.StatusOK, res)
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
