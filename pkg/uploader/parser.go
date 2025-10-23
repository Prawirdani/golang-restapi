package uploader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	httperr "github.com/prawirdani/golang-restapi/pkg/errors"
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
			return nil, httperr.PayloadTooLarge(
				fmt.Sprintf("body too large (max %d bytes)", maxBytesErr.Limit),
			)
		}

		if err == http.ErrMissingFile && !p.config.Required {
			return nil, nil // Return empty file for optional uploads
		}

		if p.config.Required {
			return nil, httperr.BadRequest(fmt.Sprintf("file '%s' is required", fieldName))
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
		return nil, err
	}

	return parsed, nil
}

// validate performs all validation checks on a file
func (p *Parser) validate(parsed *ParsedFile) error {
	// Check size
	if p.config.MaxSize > 0 && parsed.size > p.config.MaxSize {
		return httperr.PayloadTooLarge(fmt.Sprintf(
			"file size %d bytes exceeds maximum %d bytes",
			parsed.size,
			p.config.MaxSize,
		))
	}

	// Check forbidden extensions
	if len(p.config.ForbiddenExts) > 0 {
		for _, ext := range p.config.ForbiddenExts {
			if parsed.extension == strings.ToLower(ext) {
				return httperr.BadRequest(
					fmt.Sprintf("file extension %s is not allowed", parsed.extension),
				)
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
			return httperr.BadRequest(fmt.Sprintf(
				"file extension %s is not allowed. Allowed: %v",
				parsed.extension,
				p.config.AllowedExts,
			))
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
			return httperr.BadRequest(fmt.Sprintf(
				"file type %s is not allowed. Allowed: %v",
				contentType,
				p.config.AllowedMIMEs,
			))
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

// ParseMultipleFiles parses and validates multiple files from HTTP request
// func (p *Parser) ParseMultipleFiles(
// 	r *http.Request,
// 	fieldName string,
// 	maxFiles int,
// ) ([]*ParsedFile, error) {
// 	// Parse multipart form
// 	if err := r.ParseMultipartForm(p.config.MaxSize); err != nil {
// 		return nil, fmt.Errorf("files too large (max %d bytes total)", p.config.MaxSize)
// 	}
//
// 	// Get files
// 	if r.MultipartForm == nil {
// 		if p.config.Required {
// 			return nil, fmt.Errorf("no files provided")
// 		}
// 		return nil, nil
// 	}
//
// 	fileHeaders := r.MultipartForm.File[fieldName]
// 	if len(fileHeaders) == 0 {
// 		if p.config.Required {
// 			return nil, fmt.Errorf("files '%s' are required", fieldName)
// 		}
// 		return nil, nil
// 	}
//
// 	if maxFiles > 0 && len(fileHeaders) > maxFiles {
// 		return nil, fmt.Errorf("maximum %d files allowed, got %d", maxFiles, len(fileHeaders))
// 	}
//
// 	var parsedFiles []*ParsedFile
//
// 	for _, header := range fileHeaders {
// 		file, err := header.Open()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to open file %s: %w", header.Filename, err)
// 		}
//
// 		parsed := &ParsedFile{
// 			file:      file,
// 			header:    header,
// 			filename:  header.Filename,
// 			extension: strings.ToLower(filepath.Ext(header.Filename)),
// 			size:      header.Size,
// 			noFile:    false,
// 		}
//
// 		if err := p.validate(parsed); err != nil {
// 			file.Close()
// 			// Close all previously opened files
// 			for _, pf := range parsedFiles {
// 				pf.Close()
// 			}
// 			return nil, fmt.Errorf("validation failed for %s: %w", header.Filename, err)
// 		}
//
// 		parsedFiles = append(parsedFiles, parsed)
// 	}
//
// 	return parsedFiles, nil
// }
