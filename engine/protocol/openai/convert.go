package openai

import "github.com/icosmos-space/llmatrix/engine/types"

// ToChatRequest converts an engine ChatRequest into an OpenAI-shaped payload.
// Concrete wire structs will be filled in as provider support lands.
func ToChatRequest(req *types.ChatRequest) *types.ChatRequest {
	return req
}

// FromChatResponse normalizes an OpenAI response into engine types.
func FromChatResponse(resp *types.ChatResponse) *types.ChatResponse {
	return resp
}
