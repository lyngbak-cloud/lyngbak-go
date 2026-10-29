package lyngbak

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// APIError represents a non-2xx response from the Miget API.
//
// The API returns error payloads in two shapes:
//
//	{"error": "message"}
//	{"errors": ["message one", "message two"]}
//
// Both are parsed. If the body matches neither shape, the raw body is
// preserved in Body and used as a fallback in Error().
type APIError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Message is the value of the "error" (or "message") field, if present.
	Message string
	// Errors is the list of messages from the "errors" field, if present.
	Errors []string
	// Body is the raw response body.
	Body []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	var detail string
	switch {
	case e.Message != "" && len(e.Errors) > 0:
		detail = e.Message + ": " + strings.Join(e.Errors, "; ")
	case e.Message != "":
		detail = e.Message
	case len(e.Errors) > 0:
		detail = strings.Join(e.Errors, "; ")
	default:
		body := strings.TrimSpace(string(e.Body))
		if body != "" {
			const max = 200
			if len(body) > max {
				body = body[:max] + "..."
			}
			detail = body
		} else {
			detail = http.StatusText(e.StatusCode)
		}
	}
	return fmt.Sprintf("lyngbak: API error %d: %s", e.StatusCode, detail)
}

// newAPIError builds an APIError from a response status code and body.
func newAPIError(statusCode int, body []byte) *APIError {
	e := &APIError{StatusCode: statusCode, Body: body}

	var payload struct {
		Error   string          `json:"error"`
		Message string          `json:"message"`
		Errors  json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		e.Message = payload.Error
		if e.Message == "" {
			e.Message = payload.Message
		}
		if len(payload.Errors) > 0 {
			var list []string
			if err := json.Unmarshal(payload.Errors, &list); err == nil {
				e.Errors = list
			} else {
				// Some APIs return "errors" as a single string.
				var single string
				if err := json.Unmarshal(payload.Errors, &single); err == nil && single != "" {
					e.Errors = []string{single}
				}
			}
		}
	}
	return e
}

// IsNotFound reports whether err is an *APIError with HTTP status 404.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}
