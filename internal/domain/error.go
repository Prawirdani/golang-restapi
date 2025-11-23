package domain

type ErrorKind int8

const (
	ErrorKindNotFound     ErrorKind = iota // Resource not found
	ErrorKindDuplicate                     // Already exists or duplicate constraint
	ErrorKindUnauthorized                  // Not authenticated
	ErrorKindForbidden                     // No permission / action against bussiness rules
	ErrorKindValidation                    // Business rule violation
	ErrorKindTimeout                       // Operation timed out
	ErrorKindUnavailable                   // Temporarily unavailable
)

// Error is domain error type
type Error struct {
	Message string
	Kind    ErrorKind
}

func (e *Error) Error() string {
	return e.Message
}

// NewError returns new domain Error
func NewError(kind ErrorKind, message string) *Error {
	return &Error{
		Message: message,
		Kind:    kind,
	}
}

var (
	ErrNotFound     = factory(ErrorKindNotFound)
	ErrDuplicate    = factory(ErrorKindDuplicate)
	ErrUnauthorized = factory(ErrorKindUnauthorized)
	ErrForbidden    = factory(ErrorKindForbidden)
	ErrValidation   = factory(ErrorKindValidation)
	ErrTimeout      = factory(ErrorKindTimeout)
)

func factory(kind ErrorKind) func(msg string) *Error {
	return func(msg string) *Error {
		return &Error{
			Kind:    kind,
			Message: msg,
		}
	}
}
