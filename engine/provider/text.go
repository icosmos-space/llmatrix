package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// TextCompletionProvider implements OpenAI legacy /completions.
type TextCompletionProvider interface {
	Provider
	TextCompletion(ctx context.Context, req *types.TextCompletionRequest) (*types.TextCompletionResponse, error)
}
