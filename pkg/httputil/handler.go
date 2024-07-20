package httputil

import "net/http"

// customHandler is a custom handler wrapper function type.
// It is used to wrap handler function to make it easier handling errors from handler function.
type CustomHandler func(w http.ResponseWriter, r *http.Request) error
