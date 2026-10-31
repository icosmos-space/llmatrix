package engine

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/registry"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Executor runs chat and embedding calls against registered providers.
type Executor struct {
	registry *registry.Registry
}

// NewExecutor creates an Executor.
func NewExecutor(reg *registry.Registry) *Executor {
	return &Executor{registry: reg}
}

// ChatCompletion routes a chat request to the named provider.
func (e *Executor) ChatCompletion(ctx context.Context, providerName string, req *types.ChatRequest) (*types.ChatResponse, error) {
	p, ok := e.registry.Get(providerName)
	if !ok {
		return nil, fmt.Errorf("engine: provider %q not found", providerName)
	}
	chatProv, ok := p.(provider.ChatProvider)
	if !ok {
		return nil, fmt.Errorf("engine: provider %q does not support chat", providerName)
	}
	return chatProv.ChatCompletion(ctx, req)
}

// Embeddings routes an embedding request to the named provider.
func (e *Executor) Embeddings(ctx context.Context, providerName string, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	p, ok := e.registry.Get(providerName)
	if !ok {
		return nil, fmt.Errorf("engine: provider %q not found", providerName)
	}
	embProv, ok := p.(provider.EmbeddingProvider)
	if !ok {
		return nil, fmt.Errorf("engine: provider %q does not support embeddings", providerName)
	}
	return embProv.Embeddings(ctx, req)
}
