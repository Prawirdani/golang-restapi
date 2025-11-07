package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/internal/transport/http/request"
	"github.com/prawirdani/golang-restapi/internal/transport/http/request/uploader"
	"github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

var (
	// General transport layer errors
	ErrResourceNotFound = errors.New("The requested resource could not be found")
	ErrMethodNotAllowed = errors.New("The method is not allowed for the requested URL")
	ErrTooManyRequest   = errors.New("Too many request, try again later")

	// Register domain error here to have appropriate status code
	domainErrMap = map[error]int{
		// User Errors
		user.ErrEmailNotVerified: http.StatusForbidden,
		user.ErrEmailExist:       http.StatusConflict,
		user.ErrUserNotFound:     http.StatusNotFound,

		// Auth Errors
		auth.ErrWrongCredentials:           http.StatusUnauthorized,
		auth.ErrTokenExpired:               http.StatusUnauthorized,
		auth.ErrTokenInvalid:               http.StatusUnauthorized,
		auth.ErrTokenNotProvided:           http.StatusUnauthorized,
		auth.ErrSessionExpired:             http.StatusForbidden,
		auth.ErrSessionNotFound:            http.StatusUnauthorized,
		auth.ErrResetPasswordTokenInvalid:  http.StatusForbidden,
		auth.ErrResetPasswordTokenNotFound: http.StatusNotFound,
	}
)

// StatusCode returns the appropriate HTTP status code for a domain error.
// Falls back to 500 if no match is found.
func domainErrStatusCode(err error) int {
	for target, code := range domainErrMap {
		if errors.Is(err, target) {
			return code
		}
	}
	return http.StatusInternalServerError
}

type errorJsonBody struct {
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func HandleError(w http.ResponseWriter, err error) {
	var (
		validationErr *validator.ValidationError
		jsonBodyErr   *request.JsonBodyError
		uploaderErr   *uploader.ParserError
		body          errorJsonBody
		statusCode    int
	)

	switch {
	case errors.Is(err, context.Canceled):
		statusCode = http.StatusServiceUnavailable
		body.Message = "Server is busy"

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
		statusCode = domainErrStatusCode(err)
		if statusCode == http.StatusInternalServerError {
			log.Error("unknown error", err)
			body.Message = "An unexpected error occurred, try again later"
		} else {
			body.Message = err.Error()
		}
	}
	_ = response.WriteJson(w, statusCode, &body)
}
