package gemini

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/types"
)

type wirePart struct {
	Text             string            `json:"text,omitempty"`
	FunctionCall     *wireFunctionCall `json:"functionCall,omitempty"`
	FunctionResponse *wireFunctionResp `json:"functionResponse,omitempty"`
}

type wireFunctionCall struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

type wireFunctionResp struct {
	Name     string `json:"name"`
	Response any    `json:"response"`
}

type wireContent struct {
	Role  string     `json:"role,omitempty"`
	Parts []wirePart `json:"parts"`
}

type wireGenerationConfig struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
}

type wireFunctionDecl struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type wireTool struct {
	FunctionDeclarations []wireFunctionDecl `json:"functionDeclarations"`
}

type wireRequest struct {
	Contents          []wireContent         `json:"contents"`
	SystemInstruction *wireContent          `json:"systemInstruction,omitempty"`
	GenerationConfig  *wireGenerationConfig `json:"generationConfig,omitempty"`
	Tools             []wireTool            `json:"tools,omitempty"`
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
		switch m.Role {
		case types.RoleTool:
			var resp any
			if err := json.Unmarshal([]byte(m.Content), &resp); err != nil {
				resp = map[string]any{"result": m.Content}
			}
			name := m.Name
			if name == "" {
				name = "tool"
			}
			out = append(out, wireContent{
				Role: "user",
				Parts: []wirePart{{
					FunctionResponse: &wireFunctionResp{Name: name, Response: resp},
				}},
			})
		case types.RoleAssistant:
			parts := make([]wirePart, 0, 1+len(m.ToolCalls))
			if m.Content != "" {
				parts = append(parts, wirePart{Text: m.Content})
			}
			for _, tc := range m.ToolCalls {
				args := json.RawMessage(tc.Arguments)
				if !json.Valid(args) {
					b, _ := json.Marshal(tc.Arguments)
					args = b
				}
				parts = append(parts, wirePart{
					FunctionCall: &wireFunctionCall{Name: tc.Name, Args: args},
				})
			}
			if len(parts) == 0 {
				parts = append(parts, wirePart{Text: ""})
			}
			out = append(out, wireContent{Role: "model", Parts: parts})
		default:
			out = append(out, wireContent{
				Role:  "user",
				Parts: []wirePart{{Text: m.Content}},
			})
		}
	}
	return out
}

func toWireTools(tools []types.Tool) []wireTool {
	if len(tools) == 0 {
		return nil
	}
	decls := make([]wireFunctionDecl, 0, len(tools))
	for _, t := range tools {
		decls = append(decls, wireFunctionDecl{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
		})
	}
	return []wireTool{{FunctionDeclarations: decls}}
}

func partsText(c wireContent) string {
	var s string
	for _, p := range c.Parts {
		s += p.Text
	}
	return s
}

func partsToolCalls(c wireContent) []types.ToolCall {
	var out []types.ToolCall
	for i, p := range c.Parts {
		if p.FunctionCall == nil {
			continue
		}
		args := string(p.FunctionCall.Args)
		if args == "" {
			args = "{}"
		}
		out = append(out, types.ToolCall{
			ID:        fmtToolID(i, p.FunctionCall.Name),
			Type:      "function",
			Name:      p.FunctionCall.Name,
			Arguments: args,
		})
	}
	return out
}

func fmtToolID(i int, name string) string {
	return fmt.Sprintf("call_%s_%d", name, i)
}

func mapFinishReason(reason string) string {
	switch reason {
	case "STOP":
		return "stop"
	case "MAX_TOKENS":
		return "length"
	case "SAFETY":
		return "content_filter"
	case "TOOL_CALLS", "FUNCTION_CALL":
		return "tool_calls"
	default:
		return strings.ToLower(reason)
	}
}
