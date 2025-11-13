package uploader

import (
	"fmt"
	"mime/multipart"
	"path"
	"strings"
)

type ParsedFile struct {
	header      *multipart.FileHeader
	filename    string
	size        int64
	contentType string
	file        multipart.File // opened lazily
}

// NewParsedFile constructs ParsedFile metadata (does not open file yet).
func NewParsedFile(fh *multipart.FileHeader) *ParsedFile {
	if fh == nil {
		return emptyFile()
	}

	return &ParsedFile{
		header:      fh,
		filename:    fh.Filename,
		size:        fh.Size,
		contentType: fh.Header.Get("Content-Type"),
	}
}

// Open lazily opens the underlying multipart file.
// Safe to call multiple times; closes any previously opened handle.
func (pf *ParsedFile) Open() error {
	if pf.header == nil {
		return fmt.Errorf("no file available")
	}

	// Close existing file if open
	if pf.file != nil {
		if err := pf.file.Close(); err != nil {
			return fmt.Errorf("failed to close previous file: %w", err)
		}
	}

	f, err := pf.header.Open()
	if err != nil {
		pf.file = nil // Ensure consistent state
		return err
	}
	pf.file = f
	return nil
}

// Close closes the underlying file if opened.
func (pf *ParsedFile) Close() error {
	if pf.file != nil {
		err := pf.file.Close()
		pf.file = nil
		return err
	}
	return nil
}

// Read implements io.Reader.
func (pf *ParsedFile) Read(p []byte) (int, error) {
	if pf.file == nil {
		if err := pf.Open(); err != nil {
			return 0, err
		}
	}
	return pf.file.Read(p)
}

// Seek implements io.Seeker, if supported.
func (pf *ParsedFile) Seek(offset int64, whence int) (int64, error) {
	if pf.file == nil {
		if err := pf.Open(); err != nil {
			return 0, err
		}
	}
	return pf.file.Seek(offset, whence)
}

// Name implements storage.File.
func (pf *ParsedFile) Name() string {
	return pf.filename
}

// SetName implements storage.File.
func (pf *ParsedFile) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	// Strip any existing extension and add the current one
	base := strings.TrimSuffix(name, path.Ext(name))
	pf.filename = base + pf.Ext()
	return nil
}

// Ext implements storage.File.
func (pf *ParsedFile) Ext() string {
	return path.Ext(pf.filename)
}

// Size implements storage.File.
func (pf *ParsedFile) Size() int64 {
	return pf.size
}

// ContentType implements storage.File.
func (pf *ParsedFile) ContentType() string {
	return pf.contentType
}

// NoFile implements storage.File.
func (pf *ParsedFile) NoFile() bool {
	return pf.header == nil
}

// Header exposes the raw multipart header (optional).
func (pf *ParsedFile) Header() *multipart.FileHeader {
	return pf.header
}

// emptyFile produced no file for ParsedFile
func emptyFile() *ParsedFile {
	return &ParsedFile{
		filename: "",
		size:     0,
	}
}
