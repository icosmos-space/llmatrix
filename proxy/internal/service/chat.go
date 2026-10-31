package service

import (
	"context"
	"errors"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatService handles chat and embedding business logic.
type ChatService struct{}

// NewChatService creates a ChatService.
func NewChatService() *ChatService { return &ChatService{} }

// Complete runs a chat completion.
func (s *ChatService) Complete(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.New("chat service: not implemented")
}

// Embed runs an embedding request.
func (s *ChatService) Embed(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.New("embedding service: not implemented")
}
