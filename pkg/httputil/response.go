package httputil

type Response struct {
	Data    any        `json:"data"`
	Message string     `json:"message"`
	Status  int        `json:"-"`
	Error   *errorBody `json:"error"`
}

type errorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
