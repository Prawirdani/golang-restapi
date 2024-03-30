package httputil

// Response body
type response struct {
	Data any `json:"data"`
}

// Response error body
type errorResponse struct {
	Error errorBody `json:"error"`
}
type errorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
