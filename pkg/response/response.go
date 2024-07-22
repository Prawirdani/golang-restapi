package response

type Base struct {
	Data    any        `json:"data,omitempty"`
	Message string     `json:"message,omitempty"`
	Status  int        `json:"-"`
	Error   *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
