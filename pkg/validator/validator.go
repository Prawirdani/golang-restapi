package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Struct(s any) error {
	if err := v.Struct(s); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			return convertError(vErrs)
		}
		return err
	}

	return nil
}
