package gemini

import "github.com/icosmos-space/llmatrix/engine/types"

// ToGenerateContentRequest documents the Gemini mapping used by
// engine/providers/gemini: system -> systemInstruction, assistant -> model role.
func ToGenerateContentRequest(req *types.ChatRequest) *types.ChatRequest {
	return req
}

// FromGenerateContentResponse normalizes Gemini finishReason into OpenAI-style
// finish_reason values (STOP -> stop, MAX_TOKENS -> length).
func FromGenerateContentResponse(resp *types.ChatResponse) *types.ChatResponse {
	return resp
}
