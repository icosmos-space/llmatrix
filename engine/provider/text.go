package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// TextCompletionProvider implements OpenAI legacy /completions.
type TextCompletionProvider interface {
	Provider
	TextCompletion(ctx context.Context, req *types.TextCompletionRequest) (*types.TextCompletionResponse, error)
}

// TextCompletionStreamProvider streams OpenAI legacy /completions.
type TextCompletionStreamProvider interface {
	Provider
	TextCompletionStream(ctx context.Context, req *types.TextCompletionRequest) (stream.Stream, error)
}
