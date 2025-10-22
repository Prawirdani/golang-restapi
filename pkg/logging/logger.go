package logging

import (
	"github.com/prawirdani/golang-restapi/config"
)

// Logger is a custom logger interface. this interface help us to easily switch between different loggers.
type Logger interface {
	Info(cat Category, caller string, message string)
	Debug(cat Category, caller string, message string)
	Error(cat Category, caller string, message string)
	Fatal(cat Category, caller string, message string)
	// Close()
}

func NewLogger(cfg *config.Config) Logger {
	return newZeroLogger(cfg)
}
