package uploader

type ParserError struct {
	Message    string
	StatusCode int
}

func (e *ParserError) Error() string {
	return e.Message
}
