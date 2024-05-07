package httputil

import "net/http"

func GetCookie(r *http.Request, cookieName string) string {
	var value string
	if cookie, err := r.Cookie(cookieName); err == nil {
		value = cookie.Value
	}
	return value
}
