package metrics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// serve routes a request through a chi router wrapped by the instrument
// middleware, so RoutePattern() is populated when the deferred metric read runs.
func serve(m *Metrics, method, target string) {
	r := chi.NewRouter()
	r.Use(m.InstrumentHandler)
	r.Get("/users/{id}", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	req := httptest.NewRequest(method, target, nil)
	r.ServeHTTP(httptest.NewRecorder(), req)
}

func TestInstrumentHandler_UsesRoutePattern(t *testing.T) {
	m := newTestMetrics()

	serve(m, http.MethodGet, "/users/123")
	serve(m, http.MethodGet, "/users/456")

	// Both concrete paths collapse onto the "/users/{id}" template -> count 2,
	// proving the raw path is not used as a label (which would yield two series).
	got := testutil.ToFloat64(m.ReqCounter.WithLabelValues("/users/{id}", http.MethodGet, "200"))
	if got != 2 {
		t.Fatalf("template series count = %v, want 2", got)
	}

	// The concrete path must NOT exist as its own series.
	if c := testutil.ToFloat64(m.ReqCounter.WithLabelValues("/users/123", http.MethodGet, "200")); c != 0 {
		t.Errorf("raw-path series should not exist, got %v", c)
	}
}

func TestInstrumentHandler_EmptyPatternBucketed(t *testing.T) {
	m := newTestMetrics()

	// Drive the middleware with a chi RouteContext present but no matched pattern
	// (RoutePattern() == "") to exercise the "unknown" fallback branch directly.
	h := m.InstrumentHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	h.ServeHTTP(httptest.NewRecorder(), req)

	got := testutil.ToFloat64(m.ReqCounter.WithLabelValues("unknown", http.MethodGet, "404"))
	if got != 1 {
		t.Fatalf("unknown bucket count = %v, want 1", got)
	}
}
