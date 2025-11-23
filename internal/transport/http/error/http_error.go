package error

import (
	"context"
	"errors"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/domain"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

type Error struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
	status  int    `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	return e.status
}

func New(status int, msg string, details any) *Error {
	return &Error{
		Message: msg,
		status:  status,
		Errors:  details,
	}
}

func FromError(err error) *Error {
	if e, isHTTPErr := err.(*Error); isHTTPErr {
		return e
	}

	body := Error{
		status: http.StatusInternalServerError,
	}

	var e *domain.Error
	if errors.As(err, &e) {
		if status, exists := domainErrStatusCodes[e.Kind]; exists {
			body.status = status
		}
		body.Message = e.Error()
	} else {
		var (
			validationErr *validator.ValidationError
			jsonBindErr   *JSONBindError
		)
		switch {
		case errors.Is(err, context.Canceled):
			body.status = http.StatusServiceUnavailable
			body.Message = "Server is busy"

		case errors.As(err, &jsonBindErr):
			body.status = http.StatusBadRequest
			body.Message = jsonBindErr.Message

		case errors.As(err, &validationErr):
			body.status = http.StatusUnprocessableEntity
			body.Message = "Validation error"
			body.Errors = validationErr.Details

		default:
			body.Message = "An unexpected error occurred, try again later"
			log.Error("Unknown error", err)
		}
	}

	return &body
}

// domainErrStatusCodes maps domain error into http status codes
var domainErrStatusCodes = map[domain.ErrorKind]int{
	domain.ErrorKindUnauthorized: http.StatusUnauthorized,        // 401
	domain.ErrorKindForbidden:    http.StatusForbidden,           // 403
	domain.ErrorKindNotFound:     http.StatusNotFound,            // 404
	domain.ErrorKindDuplicate:    http.StatusConflict,            // 409
	domain.ErrorKindValidation:   http.StatusUnprocessableEntity, // 422
	domain.ErrorKindUnavailable:  http.StatusServiceUnavailable,  // 503
}
