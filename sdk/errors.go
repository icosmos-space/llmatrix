package sdk

import (
	"errors"
	"fmt"

	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
)

var (
	// ErrNotImplemented indicates the feature is scaffolded but not yet wired.
	ErrNotImplemented = errors.New("llmatrix/sdk: not implemented")

	// ErrInvalidConfig indicates required client configuration is missing.
	ErrInvalidConfig = errors.New("llmatrix/sdk: invalid config")

	// ErrInvalidRequest indicates the request payload is invalid.
	ErrInvalidRequest = errors.New("llmatrix/sdk: invalid request")

	// ErrUnsupported indicates the selected provider lacks a capability.
	ErrUnsupported = errors.New("llmatrix/sdk: unsupported")
)

// APIError is an OpenAI-shaped error surfaced by the SDK (LiteLLM exception mapping).
type APIError struct {
	Type       string
	Message    string
	Provider   string
	StatusCode int
	Cause      error
}

func (e *APIError) Error() string {
	if e.Provider != "" {
		return fmt.Sprintf("%s [%s]: %s", e.Type, e.Provider, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *APIError) Unwrap() error { return e.Cause }

// MapError converts engine errors into OpenAI-style SDK errors.
func MapError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := engineerrors.As(err); ok {
		api := &APIError{
			Message:    e.Message,
			Provider:   e.Provider,
			StatusCode: e.StatusCode,
			Cause:      err,
		}
		switch e.Kind {
		case engineerrors.KindUnauthorized:
			api.Type = "authentication_error"
		case engineerrors.KindRateLimited:
			api.Type = "rate_limit_error"
		case engineerrors.KindInvalidRequest:
			api.Type = "invalid_request_error"
		case engineerrors.KindNotFound:
			api.Type = "not_found_error"
		case engineerrors.KindTimeout:
			api.Type = "timeout_error"
		default:
			api.Type = "api_error"
		}
		return api
	}
	return err
}

// IsAuthenticationError reports whether err is an authentication failure.
func IsAuthenticationError(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Type == "authentication_error"
}

// IsRateLimitError reports whether err is a rate limit failure.
func IsRateLimitError(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Type == "rate_limit_error"
}
