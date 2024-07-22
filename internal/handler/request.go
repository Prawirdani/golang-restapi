package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/validator"
)

// BindValidate is a helper function to bind and validate request body
func BindValidate[T any](r *http.Request) (T, error) {
	var data T

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&data); err != nil {
		return data, err
	}

	if err := validator.Struct(data); err != nil {
		return data, err
	}

	return data, nil
}
