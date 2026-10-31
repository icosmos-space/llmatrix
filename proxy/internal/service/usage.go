package service

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// UsageService records request usage.
type UsageService struct{}

// NewUsageService creates a UsageService.
func NewUsageService() *UsageService { return &UsageService{} }

// Record stores usage for a completed request.
func (s *UsageService) Record(ctx context.Context, provider, model string, u types.Usage) error {
	_ = ctx
	_ = provider
	_ = model
	_ = u
	return nil
}
