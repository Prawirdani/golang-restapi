package model

type Response struct {
	Data    any        `json:"data"`
	Message string     `json:"message"`
	Status  int        `json:"-"`
	Error   *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
