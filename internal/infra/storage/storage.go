package storage

import (
	"context"
	"io"
	"time"
)

// Storage defines the interface for object storage operations.
type Storage interface {
	// Put stores an object and returns its URL or path
	Put(ctx context.Context, path string, reader io.Reader, contentType string) (string, error)

	// Get retrieves an object
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes an object
	Delete(ctx context.Context, path string) error

	// GetURL returns a public or signed URL for an object
	GetURL(ctx context.Context, path string, expiry time.Duration) (string, error)

	// Exists checks if an object exists
	Exists(ctx context.Context, path string) (bool, error)
}
