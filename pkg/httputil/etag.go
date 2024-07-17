package httputil

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
)

func generateEtag(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// CheckEtag checks if the request has the same Etag as the data and sets the Etag header if not.
func Etag(w http.ResponseWriter, r *http.Request, data []byte) {
	etag := generateEtag(data)
	etagHeader := r.Header.Get("If-None-Match")

	if etagHeader == etag {
		w.WriteHeader(http.StatusNotModified)
	}

	w.Header().Set("Etag", etag)
}
