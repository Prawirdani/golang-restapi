package errors

import (
	"net/http"
)

// Use these error wrappers for known errors to have precise response status codes, can be used on any abstraction layer.
// It will only set the message in ErrorResponse, if you want to provide details in the ErrorResponse you should create ApiError object manually.
var (
	BadRequest          = build(http.StatusBadRequest)          // 400
	Unauthorized        = build(http.StatusUnauthorized)        // 401
	Forbidden           = build(http.StatusForbidden)           // 403
	NotFound            = build(http.StatusNotFound)            // 404
	MethodNotAllowed    = build(http.StatusMethodNotAllowed)    // 405
	Conflict            = build(http.StatusConflict)            // 409
	Gone                = build(http.StatusGone)                // 410
	UnprocessableEntity = build(http.StatusUnprocessableEntity) // 422
	TooManyRequest      = build(http.StatusTooManyRequests)     // 429
	InternalServer      = build(http.StatusInternalServerError) // 500
	ServiceUnavailable  = build(http.StatusServiceUnavailable)  // 503
	GatewayTimeout      = build(http.StatusGatewayTimeout)      // 504
)

type ApiError struct {
	Status  int
	Message string
	Cause   interface{}
}

// Return ApiErr in string format
func (e *ApiError) Error() string {
	return e.Message
}

func build(status int) func(msg string) *ApiError {
	return func(m string) *ApiError {
		return &ApiError{
			Status:  status,
			Message: m,
		}
	}
}
