package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"
	httperr "github.com/prawirdani/golang-restapi/internal/transport/http/error"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

// Context wraps http.ResponseWriter and *http.Request with helper methods
type Context struct {
	w http.ResponseWriter
	r *http.Request
}

// Func defines the handler signature that uses the custom Context and returns an error
type Func func(c *Context) error

// JSON sends a JSON response
func (c *Context) JSON(status int, data any) error {
	// Only use ETag for successful responses (2xx)
	if status >= 200 && status < 300 {
		etag := eTag(data)
		if etag != "" {
			// Check If-None-Match header
			if match := c.Get("If-None-Match"); match == etag {
				c.w.WriteHeader(http.StatusNotModified)
				return nil
			}
			c.Set("ETag", etag)
			c.Set("Cache-Control", "private, must-revalidate")
		}
		c.Set("ETag", etag)
	}

	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(status)
	return json.NewEncoder(c.w).Encode(data)
}

// String sends a plain text response
func (c *Context) String(status int, format string, values ...any) error {
	c.w.Header().Set("Content-Type", "text/plain")
	c.w.WriteHeader(status)
	_, err := fmt.Fprintf(c.w, format, values...)
	return err
}

func (c *Context) Method() string {
	return c.r.Method
}

func (c *Context) URLPath() string {
	return c.r.URL.Path
}

// Param gets a route parameter by key from Chi's URL params
func (c *Context) Param(key string) string {
	return chi.URLParam(c.r, key)
}

// Query gets a URL query parameter
func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

// Bind unmarshals JSON request body into provided struct
func (c *Context) Bind(dst any) error {
	return json.NewDecoder(c.r.Body).Decode(dst)
}

// Status sets the HTTP status code
func (c *Context) Status(code int) *Context {
	c.w.WriteHeader(code)
	return c
}

// Set sets a response header
func (c *Context) Set(key, value string) {
	c.w.Header().Set(key, value)
}

// Get gets a request header
func (c *Context) Get(key string) string {
	return c.r.Header.Get(key)
}

// SetCookie sets cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.w, cookie)
}

// GetCookie gets cookie value
func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *Context) Context() context.Context {
	return c.r.Context()
}

func (c *Context) WithContext(ctx context.Context) *Context {
	return &Context{
		r: c.r.WithContext(ctx),
		w: c.w,
	}
}

// FormValue gets a form value by key
func (c *Context) FormValue(key string) string {
	return c.r.FormValue(key)
}

// FormFile gets a single uploaded file
func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	_, header, err := c.r.FormFile(key)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return nil, httperr.New(
				http.StatusRequestEntityTooLarge,
				"request body too large",
				map[string]int{
					"max_bytes": int(maxBytesErr.Limit),
				},
			)
		}

		return nil, err

	}

	return header, nil
}

// FormFiles gets multiple uploaded files
func (c *Context) FormFiles(key string) ([]*multipart.FileHeader, error) {
	if c.r.MultipartForm == nil {
		// TODO: Dynamic? does MaxBytesReader has impact on this?
		if err := c.r.ParseMultipartForm(32 << 20); err != nil { // 32 MB default
			return nil, err
		}
	}

	if c.r.MultipartForm != nil && c.r.MultipartForm.File != nil {
		if files := c.r.MultipartForm.File[key]; len(files) > 0 {
			return files, nil
		}
	}

	return nil, http.ErrMissingFile
}

// BindValidate is a helper function to bind and validate json request body
func (c *Context) BindValidate(dst any) error {
	if c.Get("Content-Type") != "application/json" {
		return &httperr.JSONBindError{Message: "Content-Type must be application/json"}
	}

	dec := json.NewDecoder(c.r.Body)
	if err := dec.Decode(dst); err != nil {
		return httperr.ParseJSONBindErr(err)
	}

	// If Implement JSONRequestBody interfaces
	if b, ok := dst.(JSONRequestBody); ok {
		if err := b.Sanitize(); err != nil {
			return err
		}

		return b.Validate()
	}

	// Do Manual validation
	return validator.Struct(dst)
}

// Handler converts custom HandlerFunc to standard http.HandlerFunc with structured error response handling capability
func Handler(h Func) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			w: w,
			r: r,
		}
		if err := h(c); err != nil {
			e := httperr.FromError(err)
			c.JSON(e.Status(), e)
		}
	}
}

// Wrap converts standard http.HandlerFunc to custom Func
func Wrap(h http.HandlerFunc) Func {
	return func(c *Context) error {
		h(c.w, c.r)
		return nil
	}
}

// WrapHandler converts http.Handler to custom Func
func WrapHandler(h http.Handler) Func {
	return func(c *Context) error {
		h.ServeHTTP(c.w, c.r)
		return nil
	}
}

// Middleware converts Func middleware into std middleware signature
func Middleware(mw func(Func) Func) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Handler(mw(WrapHandler(next)))
	}
}
