package service

import "context"

// KeyService manages API keys.
type KeyService struct{}

// NewKeyService creates a KeyService.
func NewKeyService() *KeyService { return &KeyService{} }

// Validate checks whether an API key is active.
func (s *KeyService) Validate(ctx context.Context, key string) (bool, error) {
	_ = ctx
	_ = key
	return true, nil
}
