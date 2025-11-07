package uploader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type ParserConfig struct {
	MaxSize       int64
	AllowedExts   []string
	AllowedMIMEs  []string
	ForbiddenExts []string
	Required      bool
	ValidateMIME  bool
}

// Parser is a struct for handling file upload. It uses as the file retriever or upload engine.
type Parser struct {
	config ParserConfig
}

func New(cfg ParserConfig) *Parser {
	return &Parser{config: cfg}
}

// ParseSingleFile parses and validates a single file from HTTP request
func (p *Parser) ParseSingleFile(r *http.Request, fieldName string) (*ParsedFile, error) {
	// Get file from form
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return nil, &ParserError{
				Message:    fmt.Sprintf("body too large: max %d bytes", maxBytesErr.Limit),
				StatusCode: http.StatusRequestEntityTooLarge,
			}
		}

		if err == http.ErrMissingFile && !p.config.Required {
			return nil, nil // Return empty file for optional uploads
		}

		if p.config.Required {
			return nil, &ParserError{
				Message:    fmt.Sprintf("file '%s' is required", fieldName),
				StatusCode: http.StatusBadRequest,
			}
		}

		return nil, err

	}

	defer file.Close()

	// Create ParsedFile
	parsed := &ParsedFile{
		file:      file,
		header:    header,
		filename:  header.Filename,
		extension: strings.ToLower(filepath.Ext(header.Filename)),
		size:      header.Size,
		noFile:    false,
	}

	// Validate file
	if err := p.validate(parsed); err != nil {
		return nil, &ParserError{Message: err.Error(), StatusCode: http.StatusBadRequest}
	}

	return parsed, nil
}

// validate performs all validation checks on a file
func (p *Parser) validate(parsed *ParsedFile) error {
	// Check size
	if p.config.MaxSize > 0 && parsed.size > p.config.MaxSize {
		return fmt.Errorf(
			"file size %d bytes exceeds maximum %d bytes",
			parsed.size,
			p.config.MaxSize,
		)
	}

	// Check forbidden extensions
	if len(p.config.ForbiddenExts) > 0 {
		for _, ext := range p.config.ForbiddenExts {
			if parsed.extension == strings.ToLower(ext) {
				return fmt.Errorf("file extension %s is not allowed", parsed.extension)
			}
		}
	}

	// Check allowed extensions
	if len(p.config.AllowedExts) > 0 {
		allowed := false
		for _, ext := range p.config.AllowedExts {
			if parsed.extension == strings.ToLower(ext) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf(
				"file extension %s is not allowed. Allowed: %v",
				parsed.extension,
				p.config.AllowedExts,
			)
		}
	}

	// Validate MIME type if required
	if p.config.ValidateMIME && len(p.config.AllowedMIMEs) > 0 {
		contentType, err := p.detectMIME(parsed.file)
		if err != nil {
			return fmt.Errorf("failed to detect file type: %w", err)
		}

		parsed.contentType = contentType

		allowed := slices.Contains(p.config.AllowedMIMEs, contentType)
		if !allowed {
			return fmt.Errorf(
				"file type %s is not allowed. Allowed: %v",
				contentType,
				p.config.AllowedMIMEs,
			)
		}

		// Reset file pointer after reading
		if _, err := parsed.Seek(0, 0); err != nil {
			return fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}

	return nil
}

// detectMIME detects actual MIME type from file content
func (p *Parser) detectMIME(file io.Reader) (string, error) {
	// Read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buffer[:n]), nil
}
