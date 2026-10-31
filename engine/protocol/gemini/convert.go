package gemini

import "github.com/icosmos-space/llmatrix/engine/types"

// ToChatRequest converts an engine ChatRequest into a Gemini-shaped payload.
func ToChatRequest(req *types.ChatRequest) *types.ChatRequest {
	return req
}

// FromChatResponse normalizes a Gemini response into engine types.
func FromChatResponse(resp *types.ChatResponse) *types.ChatResponse {
	return resp
}
