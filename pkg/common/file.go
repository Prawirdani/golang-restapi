package common

type File interface {
	// io.Reader
	Read(p []byte) (n int, err error)
	// Name returns the name of the file including the extension.
	Name() string
	// SetName sets the name of the file.
	SetName(name string) error
	// Ext returns the extension of the file.
	Ext() string
	// ContentType returns the content-type of the file
	ContentType() string
	// Return true if no file/content
	NoFile() bool
	Close() error
}
