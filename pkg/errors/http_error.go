package errors

type HttpError struct {
	Status  int
	Message string
	Cause   any
}

// Return HttpErr in string format
func (e *HttpError) Error() string {
	return e.Message
}

func build(status int) func(msg string) *HttpError {
	return func(m string) *HttpError {
		return &HttpError{
			Status:  status,
			Message: m,
		}
	}
}
