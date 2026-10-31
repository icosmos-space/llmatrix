package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// ModerationProvider classifies content for safety (OpenAI /moderations).
type ModerationProvider interface {
	Provider
	Moderation(ctx context.Context, req *types.ModerationRequest) (*types.ModerationResponse, error)
}
