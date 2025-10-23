package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
)

// ParsedFile implements the File interface
type ParsedFile struct {
	file        multipart.File
	header      *multipart.FileHeader
	filename    string
	extension   string
	size        int64
	contentType string
	noFile      bool
}

// Read implements io.Reader
func (pf *ParsedFile) Read(p []byte) (n int, err error) {
	if pf.noFile || pf.file == nil {
		return 0, io.EOF
	}
	return pf.file.Read(p)
}

// Name returns the name of the file including the extension
func (pf *ParsedFile) Name() string {
	return pf.filename
}

// SetName sets the name of the file
func (pf *ParsedFile) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	pf.filename = name + pf.Ext()
	return nil
}

// Ext returns the extension of the file
func (pf *ParsedFile) Ext() string {
	return pf.extension
}

// NoFile returns true if no file was provided
func (pf *ParsedFile) NoFile() bool {
	return pf.noFile
}

// Additional methods for ParsedFile (not part of interface but useful)

// ContentType returns the detected MIME type
func (pf *ParsedFile) ContentType() string {
	return pf.contentType
}

// Close closes the underlying file
func (pf *ParsedFile) Close() error {
	if pf.file != nil {
		return pf.file.Close()
	}
	return nil
}

// Seek seeks to a position in the file
func (pf *ParsedFile) Seek(offset int64, whence int) (int64, error) {
	if pf.noFile || pf.file == nil {
		return 0, fmt.Errorf("no file available")
	}
	if seeker, ok := pf.file.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, fmt.Errorf("file does not support seeking")
}

// Header returns the original multipart header
func (pf *ParsedFile) Header() *multipart.FileHeader {
	return pf.header
}

// EmptyFile returns an empty file (when file is optional and not provided)
func EmptyFile() *ParsedFile {
	return &ParsedFile{
		noFile:   true,
		filename: "",
	}
}
