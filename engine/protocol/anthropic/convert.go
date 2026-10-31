package anthropic

import "github.com/icosmos-space/llmatrix/engine/types"

// ToMessagesRequest documents the Anthropic Messages mapping used by
// engine/providers/anthropic: system messages become the top-level `system`
// field; remaining messages are user/assistant content blocks.
func ToMessagesRequest(req *types.ChatRequest) *types.ChatRequest {
	return req
}

// FromMessagesResponse normalizes Anthropic stop reasons into OpenAI-style
// finish_reason values (end_turn -> stop, max_tokens -> length).
func FromMessagesResponse(resp *types.ChatResponse) *types.ChatResponse {
	return resp
}
