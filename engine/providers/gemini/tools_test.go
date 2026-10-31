package gemini

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatCompletionTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body wireRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Tools) != 1 || len(body.Tools[0].FunctionDeclarations) != 1 {
			t.Fatalf("tools=%+v", body.Tools)
		}
		if body.Tools[0].FunctionDeclarations[0].Name != "get_weather" {
			t.Fatalf("name=%s", body.Tools[0].FunctionDeclarations[0].Name)
		}
		_ = json.NewEncoder(w).Encode(wireResponse{
			Candidates: []wireCandidate{{
				Content: wireContent{
					Role: "model",
					Parts: []wirePart{{
						FunctionCall: &wireFunctionCall{
							Name: "get_weather",
							Args: json.RawMessage(`{"city":"SF"}`),
						},
					}},
				},
				FinishReason: "STOP",
			}},
		})
	}))
	defer srv.Close()

	p := New("k", WithBaseURL(srv.URL))
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "gemini-2.0-flash",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "weather?"},
		},
		Tools: []types.Tool{{
			Type: "function",
			Function: types.ToolFunction{
				Name:       "get_weather",
				Parameters: map[string]any{"type": "object"},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].FinishReason != "tool_calls" {
		t.Fatalf("finish=%q", resp.Choices[0].FinishReason)
	}
	if len(resp.Choices[0].Message.ToolCalls) != 1 {
		t.Fatalf("%+v", resp.Choices[0].Message.ToolCalls)
	}
	if !strings.Contains(resp.Choices[0].Message.ToolCalls[0].Arguments, "SF") {
		t.Fatalf("args=%s", resp.Choices[0].Message.ToolCalls[0].Arguments)
	}
}
