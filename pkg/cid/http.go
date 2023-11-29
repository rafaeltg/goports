package cid

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

// HeaderKey contains the correlation id key.
const HeaderKey = "X-Request-Id"

// FromRequest returns a correlation id associated with a given http request.
// If no correlation ID is found in http request, a new random correlation id
// is returned, or an error if any.
func FromRequest(req *http.Request) (string, error) {
	cid := req.Header.Get(HeaderKey)
	if cid != "" {
		return cid, nil
	}

	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate correlation id: %w", err)
	}

	return id.String(), nil
}
