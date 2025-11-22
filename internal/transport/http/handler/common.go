package handler

import (
	"errors"
	"io"
	"net/http"

	httperr "github.com/prawirdani/golang-restapi/internal/transport/http/error"
)

const (
	// AccessTokenCookie used as access token cookie name and response body field.
	// Value based on GenerateJWT
	AccessTokenCookie = "access_token"
	// RefreshTokenCookie used as refresh token cookie name and response body field.
	// Value based on Session.ID
	RefreshTokenCookie = "refresh_token"
)

var ErrMissingAuthToken = httperr.New(
	http.StatusUnauthorized,
	"missing auth token from authorization header or http-only cookie",
	nil,
)

func isMissingFileError(err error) bool {
	if errors.Is(err, http.ErrMissingFile) || errors.Is(err, io.EOF) {
		return true
	}

	return false
}
