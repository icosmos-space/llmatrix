package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// AssistantsProvider implements OpenAI Assistants + Threads (LiteLLM assistants helpers).
type AssistantsProvider interface {
	Provider
	CreateAssistant(ctx context.Context, req *types.AssistantCreateRequest) (*types.Assistant, error)
	GetAssistant(ctx context.Context, assistantID string) (*types.Assistant, error)
	DeleteAssistant(ctx context.Context, assistantID string) error
	CreateThread(ctx context.Context, req *types.ThreadCreateRequest) (*types.Thread, error)
	GetThread(ctx context.Context, threadID string) (*types.Thread, error)
	AddMessage(ctx context.Context, threadID string, msg types.ThreadMessage) (*types.ThreadMessageObject, error)
	ListMessages(ctx context.Context, threadID string) (*types.ThreadMessageList, error)
	CreateRun(ctx context.Context, threadID string, req *types.RunCreateRequest) (*types.Run, error)
	GetRun(ctx context.Context, threadID, runID string) (*types.Run, error)
}
