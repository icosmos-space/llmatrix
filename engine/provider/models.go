package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// ModelsProvider lists available models (OpenAI GET /models).
type ModelsProvider interface {
	Provider
	ListModels(ctx context.Context) (*types.ModelList, error)
}
