package uploader

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	httperr "github.com/prawirdani/golang-restapi/internal/transport/http/error"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type ValidationRules struct {
	MaxSize      int64
	AllowedMIMEs []string
}

// Check if MIME type is in allowed list
func (r ValidationRules) isMIMEAllowed(mime string) bool {
	mimeBase := strings.Split(mime, ";")[0]
	for _, allowed := range r.AllowedMIMEs {
		if strings.EqualFold(mimeBase, allowed) {
			return true
		}
	}
	return false
}

// ValidateFile validates a parsed and prepared file based on the given rules.
// Performs checks in order of cost: existence → size → MIME type → extension verification.
func ValidateFile(ctx context.Context, f *ParsedFile, rules ValidationRules) error {
	// Existence check
	if f == nil || f.NoFile() {
		return fmt.Errorf("no file provided")
	}

	// Size check (cheap)
	if rules.MaxSize > 0 && f.Size() > rules.MaxSize {
		return httperr.New(
			http.StatusBadRequest,
			"file size exceeds maximum allowed",
			map[string]any{
				"max_bytes": rules.MaxSize,
				"received":  f.size,
			},
		)
	}

	// MIME type validation (primary security check)
	if len(rules.AllowedMIMEs) > 0 {
		claimedType := f.ContentType() // From headers - UNTRUSTED
		actualType, err := detectMIME(f)
		if err != nil {
			return fmt.Errorf("failed to detect MIME type: %w", err)
		}

		// SPECIAL CASE: Skip MIME mismatch check for Office documents
		// since they are ZIP containers and often detected as application/zip
		if isOfficeDocument(f.Ext(), actualType, claimedType) {
			// For Office documents, we trust the extension and claimed type
			// but still validate that the actual type is either ZIP or the expected Office MIME
			if !rules.isMIMEAllowed(claimedType) && !rules.isMIMEAllowed(actualType) {
				return httperr.New(
					http.StatusBadRequest,
					"file type is not allowed",
					map[string]any{"allowed_mimes": rules.AllowedMIMEs},
				)
			}
			// Office document passed special validation
			return nil
		}

		// SECURITY CHECK: Verify claimed type matches actual type
		if claimedType != "" && !mimeTypesMatch(claimedType, actualType) {
			log.WarnCtx(
				ctx,
				"MIME type mismatch, possible malicious upload",
				"claimed",
				claimedType,
				"actual",
				actualType,
			)
			return httperr.New(
				http.StatusBadRequest,
				"invalid file",
				nil,
			) // Return generic response
		}

		// Check if actual type is allowed
		if !rules.isMIMEAllowed(actualType) {
			return httperr.New(
				http.StatusBadRequest,
				"file type is not allowed",
				map[string]any{"actual": actualType, "allowed_mimes": rules.AllowedMIMEs},
			)
		}

		// EXTENSION VERIFICATION: Ensure extension matches detected MIME type
		if err := verifyExtensionMatchesMIME(f.Ext(), actualType); err != nil {
			return fmt.Errorf("file extension validation failed: %w", err)
		}
	}

	return nil
}

// verifyExtensionMatchesMIME ensures the file extension matches the detected MIME type
func verifyExtensionMatchesMIME(ext, detectedMIME string) error {
	expectedExtensions := getExpectedExtensions(detectedMIME)
	if len(expectedExtensions) == 0 {
		// No expected extensions for this MIME type, skip verification
		return nil
	}

	for _, expectedExt := range expectedExtensions {
		if strings.EqualFold(ext, expectedExt) {
			return nil // Extension matches
		}
	}

	return fmt.Errorf(
		"extension '%s' does not match detected MIME type '%s' (expected extensions: %v)",
		ext, detectedMIME, expectedExtensions,
	)
}

// getExpectedExtensions returns file extensions for a given MIME type using Go's mime package
func getExpectedExtensions(mimeType string) []string {
	mimeBase := strings.Split(mimeType, ";")[0]

	// Use Go's built-in MIME type to extension mapping
	extensions, err := mime.ExtensionsByType(mimeBase)
	if err != nil {
		// getCommonExtensionsFallback provides extensions for common types not fully covered by mime package
		// Common web formats that might not be fully covered by mime package
		fallback := map[string][]string{
			"image/webp":               {".webp"},
			"image/svg+xml":            {".svg"},
			"application/wasm":         {".wasm"},
			"font/woff":                {".woff"},
			"font/woff2":               {".woff2"},
			"application/octet-stream": {}, // Skip validation for binary files
		}

		if exts, exists := fallback[mimeBase]; exists {
			return exts
		}
		return nil
	}

	return extensions
}

// detectMIME detects the MIME type of the file content and resets the file pointer.
func detectMIME(f *ParsedFile) (string, error) {
	originalOffset, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", fmt.Errorf("failed to get current position: %w", err)
	}

	// Read first 512 bytes for MIME detection
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset position
	_, err = f.Seek(originalOffset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to reset file position: %w", err)
	}

	// For empty files, return generic MIME type
	if n == 0 {
		return "application/octet-stream", nil
	}

	return http.DetectContentType(buf[:n]), nil
}

// Helper function to compare MIME types
func mimeTypesMatch(claimed, actual string) bool {
	// Normalize by removing parameters
	claimedBase := strings.Split(claimed, ";")[0]
	actualBase := strings.Split(actual, ";")[0]

	return strings.EqualFold(strings.TrimSpace(claimedBase), strings.TrimSpace(actualBase))
}

// isOfficeDocument checks if the file is a modern Office document
// that might be detected as application/zip due to its container format
func isOfficeDocument(ext, actualType, claimedType string) bool {
	officeExtensions := map[string][]string{
		".xlsx": {
			"application/zip",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/octet-stream",
		},
		".docx": {
			"application/zip",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/octet-stream",
		},
		".pptx": {
			"application/zip",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"application/octet-stream",
		},
	}

	allowedMimes, exists := officeExtensions[strings.ToLower(ext)]
	if !exists {
		return false
	}

	// Check if the actual detected type is expected for this Office extension
	for _, allowedMime := range allowedMimes {
		if strings.Contains(actualType, allowedMime) {
			return true
		}
	}

	return false
}
