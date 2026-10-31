package errors

import "fmt"

// Kind classifies engine errors.
type Kind string

const (
	KindInvalidRequest Kind = "invalid_request"
	KindUnauthorized   Kind = "unauthorized"
	KindRateLimited    Kind = "rate_limited"
	KindProvider       Kind = "provider"
	KindNotFound       Kind = "not_found"
	KindInternal       Kind = "internal"
)

// Error is a structured engine error.
type Error struct {
	Kind     Kind
	Provider string
	Message  string
	Cause    error
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
