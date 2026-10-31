package sdk

import "errors"

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
