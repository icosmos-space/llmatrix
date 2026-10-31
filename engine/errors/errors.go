package errors

import (
	"errors"
	"fmt"
)

// Kind classifies engine errors.
type Kind string

const (
	KindInvalidRequest Kind = "invalid_request"
	KindUnauthorized   Kind = "unauthorized"
	KindRateLimited    Kind = "rate_limited"
	KindProvider       Kind = "provider"
	KindNotFound       Kind = "not_found"
	KindInternal       Kind = "internal"
	KindTimeout        Kind = "timeout"
)

// Error is a structured engine error.
type Error struct {
	Kind       Kind
	Provider   string
	Message    string
	StatusCode int
	Cause      error
}

func (e *Error) Error() string {
	if e.Provider != "" {
		return fmt.Sprintf("engine/%s [%s]: %s", e.Kind, e.Provider, e.Message)
	}
	return fmt.Sprintf("engine/%s: %s", e.Kind, e.Message)
}

func (e *Error) Unwrap() error { return e.Cause }

// Wrap builds a structured Error.
func Wrap(kind Kind, provider, message string, cause error) *Error {
	return &Error{Kind: kind, Provider: provider, Message: message, Cause: cause}
}

// WrapStatus builds a structured Error with an HTTP status.
func WrapStatus(kind Kind, provider, message string, status int, cause error) *Error {
	return &Error{Kind: kind, Provider: provider, Message: message, StatusCode: status, Cause: cause}
}

// As extracts a structured *Error from err.
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// IsRetryable reports whether the error is safe to retry against the same provider.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	e, ok := As(err)
	if !ok {
		// network / context deadline often wrapped without Kind
		return true
	}
	switch e.Kind {
	case KindRateLimited, KindTimeout, KindInternal, KindProvider:
		if e.StatusCode == 400 || e.StatusCode == 401 || e.StatusCode == 403 || e.StatusCode == 404 {
			return false
		}
		return true
	default:
		return false
	}
}

// IsFallbackable reports whether another model/provider candidate should be tried.
// Invalid requests are not fallbackable; auth / rate-limit / provider errors are.
func IsFallbackable(err error) bool {
	if err == nil {
		return false
	}
	e, ok := As(err)
	if !ok {
		return true
	}
	return e.Kind != KindInvalidRequest
}
