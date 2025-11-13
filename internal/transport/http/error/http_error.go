package error

import (
	"context"
	"errors"
	"net/http"

	"github.com/prawirdani/golang-restapi/pkg/errorsx"
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

	// Prioritize the CategorizedError
	var e errorsx.CategorizedError
	if errors.As(err, &e) {
		if status, exists := catdErrStatusCodes[e.Category()]; exists {
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
			body.Message = "server is busy"

		case errors.As(err, &jsonBindErr):
			body.status = http.StatusBadRequest
			body.Message = jsonBindErr.Message

		case errors.As(err, &validationErr):
			body.status = http.StatusUnprocessableEntity
			body.Message = "Validation error"
			body.Errors = validationErr.Details

		default:
			body.Message = "an unexpected error occurred, try again later"
			log.Error("Unknown error", err)
		}
	}

	return &body
}

// catdErrStatusCodes maps CategroziedError into http status codes
var catdErrStatusCodes = map[errorsx.Category]int{
	errorsx.CategoryValidation:        http.StatusBadRequest,         // 400
	errorsx.CategoryFormat:            http.StatusBadRequest,         // 400
	errorsx.CategoryUnauthorized:      http.StatusUnauthorized,       // 401
	errorsx.CategoryForbidden:         http.StatusForbidden,          // 403
	errorsx.CategoryNotExists:         http.StatusNotFound,           // 404
	errorsx.CategoryTimeout:           http.StatusRequestTimeout,     // 408
	errorsx.CategoryDuplicate:         http.StatusConflict,           // 409
	errorsx.CategoryDependency:        http.StatusServiceUnavailable, // 503
	errorsx.CategoryUnavailable:       http.StatusServiceUnavailable, // 503
	errorsx.CategoryDependencyTimeout: http.StatusGatewayTimeout,     // 504
}
