package anthropic

import "github.com/icosmos-space/llmatrix/engine/types"

// ToChatRequest converts an engine ChatRequest into an Anthropic-shaped payload.
func ToChatRequest(req *types.ChatRequest) *types.ChatRequest {
	return req
}

// FromChatResponse normalizes an Anthropic response into engine types.
func FromChatResponse(resp *types.ChatResponse) *types.ChatResponse {
	return resp
}
