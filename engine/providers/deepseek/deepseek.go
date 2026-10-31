package deepseek

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Provider implements the DeepSeek API (OpenAI-compatible).
type Provider struct {
	APIKey  string
	BaseURL string
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	return provider.Info{
		Name:        "deepseek",
		DisplayName: "DeepSeek",
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true, Embeddings: false,
		},
	}
}

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.Wrap(errors.KindInternal, "deepseek", "not implemented", nil)
}
