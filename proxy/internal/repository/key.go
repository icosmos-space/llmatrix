package repository

import "context"

// Key is an API key record.
type Key struct {
	ID     string
	UserID string
	Secret string
	Active bool
}

// KeyRepository persists API keys.
type KeyRepository interface {
	GetBySecret(ctx context.Context, secret string) (*Key, error)
	Create(ctx context.Context, key *Key) error
}
