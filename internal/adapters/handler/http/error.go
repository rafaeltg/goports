package http

import "errors"

type (
	ErrorResponse struct {
		Error ErrorData `json:"error"`
	}

	ErrorData struct {
		Message string `json:"message"`
	}
)

var errBadRequest = errors.New("failed to read request body")
