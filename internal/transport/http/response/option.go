package response

type Option func(*jsonBody)

// WithData is an option to set data field in ResponseBody for Send function
// Pass pointer to a value if the data is a single struct with custom implementation of MarshalJSON.
// No need to pass pointer if its a slices, since JSON Marshaller will referencing each item pointer automatically
func WithData(v any) Option {
	return func(r *jsonBody) {
		r.Data = v
	}
}

// WithMessage is an option to set message field in ResponseBody for Send function
func WithMessage(msg string) Option {
	return func(r *jsonBody) {
		r.Message = msg
	}
}

// WithStatus is an option to override default OK status in ResponseBody for Send function
func WithStatus(status int) Option {
	return func(r *jsonBody) {
		r.status = status
	}
}
