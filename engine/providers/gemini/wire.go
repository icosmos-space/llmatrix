package gemini

import (
	"strings"

	"github.com/icosmos-space/llmatrix/engine/types"
)

type wirePart struct {
	Text string `json:"text,omitempty"`
}

type wireContent struct {
	Role  string     `json:"role,omitempty"`
	Parts []wirePart `json:"parts"`
}

type wireGenerationConfig struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
}

type wireRequest struct {
	Contents          []wireContent         `json:"contents"`
	SystemInstruction *wireContent          `json:"systemInstruction,omitempty"`
	GenerationConfig  *wireGenerationConfig `json:"generationConfig,omitempty"`
}

type wireCandidate struct {
	Content      wireContent `json:"content"`
	FinishReason string      `json:"finishReason"`
}

type wireUsage struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

type wireResponse struct {
	Candidates    []wireCandidate `json:"candidates"`
	UsageMetadata *wireUsage      `json:"usageMetadata,omitempty"`
	Error         *wireError      `json:"error,omitempty"`
}

type wireError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type wireEmbedRequest struct {
	Model   string      `json:"model"`
	Content wireContent `json:"content"`
}

type wireEmbedResponse struct {
	Embedding *struct {
		Values []float64 `json:"values"`
	} `json:"embedding"`
	Error *wireError `json:"error,omitempty"`
}

func splitSystem(msgs []types.Message) (system string, rest []types.Message) {
	rest = make([]types.Message, 0, len(msgs))
	for _, m := range msgs {
		if m.Role == types.RoleSystem {
			if system != "" {
				system += "\n"
			}
			system += m.Content
			continue
		}
		rest = append(rest, m)
	}
	return system, rest
}

func toWireContents(msgs []types.Message) []wireContent {
	out := make([]wireContent, 0, len(msgs))
	for _, m := range msgs {
		role := "user"
		if m.Role == types.RoleAssistant {
			role = "model"
		}
		out = append(out, wireContent{
			Role:  role,
			Parts: []wirePart{{Text: m.Content}},
		})
	}
	return out
}

func partsText(c wireContent) string {
	var s string
	for _, p := range c.Parts {
		s += p.Text
	}
	return s
}

func mapFinishReason(reason string) string {
	switch reason {
	case "STOP":
		return "stop"
	case "MAX_TOKENS":
		return "length"
	case "SAFETY":
		return "content_filter"
	default:
		return strings.ToLower(reason)
	}
}
