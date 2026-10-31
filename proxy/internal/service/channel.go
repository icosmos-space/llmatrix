package service

import "context"

// ChannelService manages upstream provider channels.
type ChannelService struct{}

// NewChannelService creates a ChannelService.
func NewChannelService() *ChannelService { return &ChannelService{} }

// ListModels returns available model identifiers.
func (s *ChannelService) ListModels(ctx context.Context) ([]string, error) {
	_ = ctx
	return []string{}, nil
}
