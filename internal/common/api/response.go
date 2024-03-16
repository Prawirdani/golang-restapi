package api

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
}
