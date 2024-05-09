package httputil

type Response struct {
	Data    any        `json:"data,omitempty"`
	Message *string    `json:"message,omitempty"`
	Status  int        `json:"-"`
	Error   *errorBody `json:"error,omitempty"`
}

type errorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
