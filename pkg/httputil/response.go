package httputil

// Response body
type response struct {
	Data any `json:"data"`
}

// Response error body
type errorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
