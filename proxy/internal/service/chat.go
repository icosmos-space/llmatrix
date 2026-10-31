package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/registry"
	"github.com/icosmos-space/llmatrix/engine/retry"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatService handles chat and embedding via Engine providers.
type ChatService struct {
	HTTPClient *http.Client
	Retry      retry.Policy
	// APIKeys optionally overrides per-provider keys (else env vars are used).
	APIKeys map[string]string
	// BaseURLs optionally overrides per-provider base URLs.
	BaseURLs map[string]string
}

// NewChatService creates a ChatService with default retry policy.
func NewChatService() *ChatService {
	return &ChatService{Retry: retry.DefaultPolicy()}
}

// Complete runs a chat completion.
func (s *ChatService) Complete(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("chat service: request is nil")
	}
	if req.Model == "" {
		return nil, fmt.Errorf("chat service: model is required")
	}
	ref := model.Resolve(req.Model, "")
	p, err := s.providerFor(ref)
	if err != nil {
		return nil, err
	}
	chatProv, ok := p.(provider.ChatProvider)
	if !ok {
		return nil, fmt.Errorf("chat service: provider %q does not support chat", ref.Provider)
	}
	engineReq := *req
	engineReq.Model = ref.Model
	engineReq.Stream = false

	var out *types.ChatResponse
	err = retry.Do(ctx, s.retryPolicy(), func(ctx context.Context) error {
		resp, callErr := chatProv.ChatCompletion(ctx, &engineReq)
		if callErr != nil {
			return callErr
		}
		out = resp
		return nil
	})
	return out, err
}

// Embed runs an embedding request.
func (s *ChatService) Embed(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("embedding service: request is nil")
	}
	if req.Model == "" {
		return nil, fmt.Errorf("embedding service: model is required")
	}
	ref := model.Resolve(req.Model, "")
	p, err := s.providerFor(ref)
	if err != nil {
		return nil, err
	}
	embProv, ok := p.(provider.EmbeddingProvider)
	if !ok {
		return nil, fmt.Errorf("embedding service: provider %q does not support embeddings", ref.Provider)
	}
	engineReq := *req
	engineReq.Model = ref.Model

	var out *types.EmbeddingResponse
	err = retry.Do(ctx, s.retryPolicy(), func(ctx context.Context) error {
		resp, callErr := embProv.Embeddings(ctx, &engineReq)
		if callErr != nil {
			return callErr
		}
		out = resp
		return nil
	})
	return out, err
}

func (s *ChatService) providerFor(ref model.Ref) (provider.Provider, error) {
	cfg := registry.ProviderConfig{
		Name:       ref.Provider,
		HTTPClient: s.HTTPClient,
	}
	if s.APIKeys != nil {
		cfg.APIKey = s.APIKeys[ref.Provider]
	}
	if s.BaseURLs != nil {
		cfg.BaseURL = s.BaseURLs[ref.Provider]
	}
	return registry.NewProvider(cfg)
}

func (s *ChatService) retryPolicy() retry.Policy {
	if s.Retry.MaxAttempts > 0 {
		return s.Retry
	}
	return retry.DefaultPolicy()
}
