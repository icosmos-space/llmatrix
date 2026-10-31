package sdk

import (
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func dropUnsupportedChatParams(req *chat.Request, ref model.Ref) *chat.Request {
	if req == nil {
		return nil
	}
	allowed := map[string]bool{}
	for _, p := range model.GetSupportedOpenAIParams(ref.Provider+"/"+ref.Model, ref.Provider) {
		allowed[p] = true
	}
	out := *req
	if !allowed["temperature"] {
		out.Temperature = nil
	}
	if !allowed["top_p"] {
		out.TopP = nil
	}
	if !allowed["n"] {
		out.N = nil
	}
	if !allowed["stop"] {
		out.Stop = nil
	}
	if !allowed["max_tokens"] {
		out.MaxTokens = nil
	}
	if !allowed["max_completion_tokens"] {
		out.MaxCompletionTokens = nil
	}
	if !allowed["presence_penalty"] {
		out.PresencePenalty = nil
	}
	if !allowed["frequency_penalty"] {
		out.FrequencyPenalty = nil
	}
	if !allowed["logit_bias"] {
		out.LogitBias = nil
	}
	if !allowed["user"] {
		out.User = ""
	}
	if !allowed["seed"] {
		out.Seed = nil
	}
	if !allowed["response_format"] {
		out.ResponseFormat = nil
	}
	if !allowed["stream_options"] {
		out.StreamOptions = nil
	}
	if !allowed["tools"] {
		out.Tools = nil
	}
	if !allowed["tool_choice"] {
		out.ToolChoice = nil
	}
	if !allowed["parallel_tool_calls"] {
		out.ParallelToolCalls = nil
	}
	if !allowed["logprobs"] {
		out.Logprobs = nil
	}
	if !allowed["top_logprobs"] {
		out.TopLogprobs = nil
	}
	if !allowed["reasoning_effort"] {
		out.ReasoningEffort = ""
	}
	if !allowed["prompt_cache_key"] {
		out.PromptCacheKey = ""
	}
	if !allowed["thinking"] {
		out.Thinking = nil
	}
	return &out
}
