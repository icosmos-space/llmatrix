package gemini

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Provider implements the Google Gemini API.
type Provider struct {
	APIKey  string
	BaseURL string
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	return provider.Info{
		Name:        "gemini",
		DisplayName: "Google Gemini",
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true, Embeddings: true,
		},
	}
}

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.Wrap(errors.KindInternal, "gemini", "not implemented", nil)
}

// Embeddings implements provider.EmbeddingProvider.
func (p *Provider) Embeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.Wrap(errors.KindInternal, "gemini", "not implemented", nil)
}
