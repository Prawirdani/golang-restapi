package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/transport/http/request"
	"github.com/prawirdani/golang-restapi/internal/transport/http/request/uploader"
	"github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/errorsx"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

var (
	// General transport layer errors
	ErrResourceNotFound = errors.New("the requested resource could not be found")
	ErrMethodNotAllowed = errors.New("the method is not allowed for the requested url")
	ErrTooManyRequest   = errors.New("too many request, try again later")
)

type errorJsonBody struct {
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func HandleError(w http.ResponseWriter, err error) {
	var (
		body       errorJsonBody
		statusCode = http.StatusInternalServerError // Default
	)

	// Prioritize the CategorizedError
	var e errorsx.CategorizedError
	if errors.As(err, &e) {
		if status, exists := catdErrStatusCodes[e.Category()]; exists {
			statusCode = status
		}
		body.Message = e.Error()
	} else {
		var (
			validationErr *validator.ValidationError
			jsonBodyErr   *request.JsonBodyError
			uploaderErr   *uploader.ParserError
		)
		switch {
		case errors.Is(err, context.Canceled):
			statusCode = http.StatusServiceUnavailable
			body.Message = "server is busy"

		case errors.Is(err, ErrMethodNotAllowed):
			statusCode = http.StatusMethodNotAllowed
			body.Message = err.Error()

		case errors.Is(err, ErrResourceNotFound):
			statusCode = http.StatusNotFound
			body.Message = err.Error()

		case errors.Is(err, ErrTooManyRequest):
			statusCode = http.StatusTooManyRequests
			body.Message = err.Error()

		case errors.As(err, &jsonBodyErr):
			statusCode = http.StatusBadRequest
			body.Message = jsonBodyErr.Message

		case errors.As(err, &uploaderErr):
			statusCode = uploaderErr.StatusCode
			body.Message = uploaderErr.Message

		case errors.As(err, &validationErr):
			statusCode = http.StatusUnprocessableEntity
			body.Message = "Validation error"
			body.Errors = validationErr.Details

		default:
			body.Message = "an unexpected error occurred, try again later"
		}
	}

	_ = response.WriteJson(w, statusCode, &body)
}

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
