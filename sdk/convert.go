package sdk

import (
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
	"github.com/icosmos-space/llmatrix/sdk/tools"
)

func toEngineChatRequest(req *chat.Request, model string) *types.ChatRequest {
	out := &types.ChatRequest{
		Model:       model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      req.Stream,
		Messages:    make([]types.Message, 0, len(req.Messages)),
	}
	for _, m := range req.Messages {
		msg := types.Message{
			Role:       types.Role(m.Role),
			Content:    m.Content,
			Name:       m.Name,
			ToolCallID: m.ToolCallID,
		}
		for _, tc := range m.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, types.ToolCall{
				ID:        tc.ID,
				Type:      tc.Type,
				Name:      tc.Name,
				Arguments: tc.Arguments,
			})
		}
		out.Messages = append(out.Messages, msg)
	}
	for _, t := range req.Tools {
		out.Tools = append(out.Tools, types.Tool{
			Type: t.Type,
			Function: types.ToolFunction{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		})
	}
	return out
}

func fromEngineChatResponse(resp *types.ChatResponse) *chat.Response {
	if resp == nil {
		return nil
	}
	out := &chat.Response{
		ID:       resp.ID,
		Model:    resp.Model,
		Provider: resp.Provider,
		Usage: chat.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Choices: make([]chat.Choice, 0, len(resp.Choices)),
	}
	for _, c := range resp.Choices {
		msg := chat.Message{
			Role:       string(c.Message.Role),
			Content:    c.Message.Content,
			Name:       c.Message.Name,
			ToolCallID: c.Message.ToolCallID,
		}
		for _, tc := range c.Message.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, tools.Call{
				ID:        tc.ID,
				Type:      tc.Type,
				Name:      tc.Name,
				Arguments: tc.Arguments,
			})
		}
		out.Choices = append(out.Choices, chat.Choice{
			Index:        c.Index,
			Message:      msg,
			FinishReason: c.FinishReason,
		})
	}
	return out
}

func toEngineEmbeddingRequest(req *embedding.Request, model string) *types.EmbeddingRequest {
	return &types.EmbeddingRequest{
		Model: model,
		Input: req.Input,
	}
}

func fromEngineEmbeddingResponse(resp *types.EmbeddingResponse) *embedding.Response {
	if resp == nil {
		return nil
	}
	out := &embedding.Response{
		Model:    resp.Model,
		Provider: resp.Provider,
		Usage: embedding.Usage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
		Data: make([]embedding.Item, 0, len(resp.Data)),
	}
	for _, d := range resp.Data {
		out.Data = append(out.Data, embedding.Item{
			Index:     d.Index,
			Embedding: d.Embedding,
		})
	}
	return out
}
