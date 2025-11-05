package response

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	stderrs "errors"

	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

type Body struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"-"`
	Errors  any    `json:"errors,omitempty"`
}

// Send is a function to send json response to client
// It uses option pattern to accepts multiple options to customize the response
func Send(w http.ResponseWriter, r *http.Request, opts ...Option) error {
	res := Body{
		Status: http.StatusOK, // Default
	}
	for _, opt := range opts {
		opt(&res)
	}

	etag := eTag(res)
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	w.Header().Set("ETag", etag)
	return writeJSON(w, &res)
}

func HandleError(w http.ResponseWriter, err error) {
	var (
		body          Body
		httpErr       *errors.HttpError
		validationErr *validator.ValidationError
	)

	switch {
	case stderrs.Is(err, context.Canceled):
		body = Body{
			Status:  http.StatusServiceUnavailable,
			Message: "Server is busy",
		}

	case stderrs.As(err, &httpErr):
		body = Body{
			Status:  httpErr.Status,
			Message: httpErr.Message,
			Errors:  httpErr.Cause,
		}

	case stderrs.As(err, &validationErr):
		body = Body{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation error",
			Errors:  validationErr.Details,
		}

	default:
		log.Error("unknown error", err)
		body = Body{
			Status:  http.StatusInternalServerError,
			Message: "An unexpected error occurred, try again later",
		}
	}
	_ = writeJSON(w, &body)
}

func eTag(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "error"
	}
	// Compute the SHA-256 hash of the JSON representation
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf(`"%x"`, hash)
}
