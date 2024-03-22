package utils

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	// Create validator instance on init
	Validate = validator.New()
}
