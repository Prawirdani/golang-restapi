package utils

import (
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	// Create validator instance on init
	v = validator.New()
}

func ValidateStruct(s interface{}) error {
	return v.Struct(s)
}
