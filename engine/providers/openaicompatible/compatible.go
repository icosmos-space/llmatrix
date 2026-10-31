package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Provider talks to any OpenAI-compatible HTTP API.
type Provider struct {
	Name    string
	APIKey  string
	BaseURL string
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	name := p.Name
	if name == "" {
		name = "openai-compatible"
	}
	return provider.Info{
		Name:        name,
		DisplayName: name,
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true, Embeddings: true,
		},
	}
}

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "not implemented", nil)
}

// Embeddings implements provider.EmbeddingProvider.
func (p *Provider) Embeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "not implemented", nil)
}
