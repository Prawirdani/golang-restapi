package handler

import (
	"errors"
	"io"
	"net/http"
)

const BindValidateWarnLog = "Failed to bind & validate json req body"

func isMissingFileError(err error) bool {
	if errors.Is(err, http.ErrMissingFile) || errors.Is(err, io.EOF) {
		return true
	}

	return false
}
