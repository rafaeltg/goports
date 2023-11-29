package http

import (
	"encoding/json"
	"net/http"
)

type (
	response struct {
		code int
		body any
	}

	responseOption func(*response)
)

func withStatusCode(code int) responseOption {
	return func(r *response) {
		r.code = code
	}
}

func withBody(body any) responseOption {
	return func(r *response) {
		r.body = body
	}
}

func withError(err error) responseOption {
	return func(r *response) {
		r.body = ErrorResponse{
			Error: ErrorData{
				Message: err.Error(),
			},
		}
	}
}

func writeResponse(w http.ResponseWriter, opts ...responseOption) {
	r := response{}
	for _, opt := range opts {
		opt(&r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.code)

	if r.body == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(r.body)
}
