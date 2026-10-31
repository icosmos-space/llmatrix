package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatProvider performs chat completions.
type ChatProvider interface {
	Provider
	ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error)
}
