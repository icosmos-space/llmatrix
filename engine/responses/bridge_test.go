package responses

import (
	"context"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
)

type fakeChat struct{}

func (fakeChat) Info() provider.Info {
	return provider.Info{Name: "fake"}
}

func (fakeChat) ChatCompletion(_ context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	return &types.ChatResponse{
		ID: "chat_1", Model: req.Model, Provider: "fake",
		Choices: []types.ChatChoice{{
			Message: types.Message{Role: types.RoleAssistant, Content: "bridged"},
		}},
		Usage: types.Usage{PromptTokens: 1, CompletionTokens: 1, TotalTokens: 2},
	}, nil
}

func TestBridgeViaChat(t *testing.T) {
	resp, err := Create(context.Background(), fakeChat{}, &types.ResponsesRequest{
		Model: "m",
		Input: "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.OutputText() != "bridged" {
		t.Fatalf("got %q", resp.OutputText())
	}
	if resp.Object != "response" {
		t.Fatalf("object=%q", resp.Object)
	}
}
