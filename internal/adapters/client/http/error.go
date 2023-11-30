package http

type (
	// ApiError represents the error data.
	ApiError struct {
		Message string `json:"message"`
	}

	// ApiErrorResponse represents the error response payload.
	ApiErrorResponse struct {
		Err ApiError `json:"error"`
	}
)

func (e ApiError) Error() string {
	return e.Message
}

func (r ApiErrorResponse) Error() string {
	return r.Err.Error()
}
