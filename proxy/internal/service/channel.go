package service

import "context"

// ChannelService manages upstream provider channels.
type ChannelService struct {
	Models []string
}

// NewChannelService creates a ChannelService with built-in model examples.
func NewChannelService() *ChannelService {
	return &ChannelService{
		Models: []string{
			"openai/gpt-4o-mini",
			"deepseek/deepseek-chat",
			"ollama/llama3",
			"anthropic/claude-sonnet-4-5",
			"gemini/gemini-2.0-flash",
		},
	}
}

// ListModels returns available model identifiers.
func (s *ChannelService) ListModels(ctx context.Context) ([]string, error) {
	_ = ctx
	out := make([]string, len(s.Models))
	copy(out, s.Models)
	return out, nil
}
