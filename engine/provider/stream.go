package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatStreamProvider performs streaming chat completions.
type ChatStreamProvider interface {
	Provider
	ChatCompletionStream(ctx context.Context, req *types.ChatRequest) (stream.Stream, error)
}
