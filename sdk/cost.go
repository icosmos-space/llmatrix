package sdk

import (
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/engine/usage"
	"github.com/icosmos-space/llmatrix/sdk/chat"
)

// TokenCounter estimates tokens (LiteLLM token_counter; approximate, not tiktoken).
func TokenCounter(text string) int {
	return usage.TokenCounter{}.CountText(text)
}

// TokenCounterMessages estimates prompt tokens for chat messages.
func TokenCounterMessages(msgs []chat.Message) int {
	engineMsgs := make([]types.Message, 0, len(msgs))
	for _, m := range msgs {
		msg := types.Message{
			Role: types.Role(m.Role), Content: m.Content, Name: m.Name, ToolCallID: m.ToolCallID,
		}
		for _, part := range m.Parts {
			cp := types.ContentPart{Type: part.Type, Text: part.Text}
			msg.Parts = append(msg.Parts, cp)
		}
		for _, tc := range m.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, types.ToolCall{
				ID: tc.ID, Type: tc.Type, Name: tc.Name, Arguments: tc.Arguments,
			})
		}
		engineMsgs = append(engineMsgs, msg)
	}
	return usage.TokenCounter{}.CountMessages(engineMsgs)
}

// CompletionCost estimates USD cost for a completion (LiteLLM completion_cost).
// modelRef is provider/model or an alias; usage comes from the response.
func CompletionCost(modelRef string, promptTokens, completionTokens int) float64 {
	ref := model.ResolveAlias(model.DefaultAliases(), modelRef, "")
	key := ref.Provider + "/" + ref.Model
	return usage.CostUSD(key, types.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
	}, nil)
}

// CompletionCostFromUsage is CompletionCost using a chat.Usage value.
func CompletionCostFromUsage(modelRef string, u chat.Usage) float64 {
	return CompletionCost(modelRef, u.PromptTokens, u.CompletionTokens)
}
