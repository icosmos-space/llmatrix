package auth

import "errors"

// ErrUnauthorized indicates authentication failure.
var ErrUnauthorized = errors.New("auth: unauthorized")

// APIKey is a persisted API key record.
type APIKey struct {
	ID     string
	UserID string
	Key    string
	Active bool
}
