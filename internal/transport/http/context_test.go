package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newCtx(method string, header http.Header) (*Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, "/", nil)
	if header != nil {
		req.Header = header
	}
	rec := httptest.NewRecorder()
	return &Context{w: rec, r: req}, rec
}

func TestJSON_BodyShape(t *testing.T) {
	c, rec := newCtx(http.MethodPost, nil)

	if err := c.JSON(&Body{Data: map[string]any{"id": 1}, Message: "ok"}); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q", ct)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := got["data"]; !ok {
		t.Errorf("missing data field: %v", got)
	}
	if got["message"] != "ok" {
		t.Errorf("message = %v, want ok", got["message"])
	}
}

func TestJSON_EmptyMessageIsNull(t *testing.T) {
	c, rec := newCtx(http.MethodPost, nil)

	if err := c.JSON(&Body{Data: 1}); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := got["message"]; !ok {
		t.Fatal("message field must be present")
	}
	if string(got["message"]) != "null" {
		t.Errorf("message = %s, want null", got["message"])
	}
}

func TestJSON_ETagOnlyForReads(t *testing.T) {
	tests := []struct {
		method   string
		wantETag bool
	}{
		{http.MethodGet, true},
		{http.MethodHead, true},
		{http.MethodPost, false},
		{http.MethodPut, false},
		{http.MethodDelete, false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			c, rec := newCtx(tt.method, nil)
			if err := c.JSON(&Body{Data: "x"}); err != nil {
				t.Fatalf("JSON: %v", err)
			}
			gotETag := rec.Header().Get("ETag") != ""
			if gotETag != tt.wantETag {
				t.Errorf("ETag present = %v, want %v", gotETag, tt.wantETag)
			}
		})
	}
}

func TestJSON_NoETagOnErrorStatus(t *testing.T) {
	c, rec := newCtx(http.MethodGet, nil)
	if err := c.JSON(&Body{Data: "x"}); err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if etag := rec.Header().Get("ETag"); etag != "" {
		t.Errorf("ETag = %q, want empty on 5xx", etag)
	}
}

func TestJSON_IfNoneMatchReturns304(t *testing.T) {
	// First GET to capture the ETag.
	c1, rec1 := newCtx(http.MethodGet, nil)
	if err := c1.JSON(&Body{Data: "cacheable"}); err != nil {
		t.Fatalf("JSON: %v", err)
	}
	etag := rec1.Header().Get("ETag")
	if etag == "" {
		t.Fatal("expected ETag on GET")
	}

	// Second GET with matching If-None-Match -> 304, empty body.
	h := http.Header{}
	h.Set("If-None-Match", etag)
	c2, rec2 := newCtx(http.MethodGet, h)
	if err := c2.JSON(&Body{Data: "cacheable"}); err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if rec2.Code != http.StatusNotModified {
		t.Errorf("status = %d, want 304", rec2.Code)
	}
	if rec2.Body.Len() != 0 {
		t.Errorf("body should be empty on 304, got %q", rec2.Body.String())
	}
}
