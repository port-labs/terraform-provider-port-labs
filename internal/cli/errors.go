package cli

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError is returned when the Port API answers with a non-OK response. It keeps the HTTP status code so callers can
// branch on it instead of matching on the response text.
type APIError struct {
	Op         string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s, got: %s", e.Op, e.Body)
}

// IsNotFound reports whether err is (or wraps) an APIError with HTTP status 404.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}
