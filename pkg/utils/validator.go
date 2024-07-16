package utils

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func Validate(s interface{}) error {
	return v.Struct(s)
}
