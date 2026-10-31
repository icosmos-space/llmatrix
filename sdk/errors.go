package sdk

import "errors"

var (
	// ErrNotImplemented indicates the feature is scaffolded but not yet wired.
	ErrNotImplemented = errors.New("llmatrix/sdk: not implemented")

	// ErrInvalidConfig indicates required client configuration is missing.
	ErrInvalidConfig = errors.New("llmatrix/sdk: invalid config")
)
