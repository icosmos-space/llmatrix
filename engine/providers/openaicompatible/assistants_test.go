package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestAssistantsFlow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
			t.Fatalf("beta=%q", r.Header.Get("OpenAI-Beta"))
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/assistants":
			_ = json.NewEncoder(w).Encode(types.Assistant{ID: "asst_1", Model: "gpt-4o", Name: "Math"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/threads":
			_ = json.NewEncoder(w).Encode(types.Thread{ID: "thread_1"})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/messages"):
			_ = json.NewEncoder(w).Encode(types.ThreadMessageObject{ID: "msg_1", Role: "user", ThreadID: "thread_1"})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/runs"):
			_ = json.NewEncoder(w).Encode(types.Run{ID: "run_1", Status: "queued", ThreadID: "thread_1", AssistantID: "asst_1"})
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/runs/"):
			_ = json.NewEncoder(w).Encode(types.Run{ID: "run_1", Status: "completed"})
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	asst, err := p.CreateAssistant(context.Background(), &types.AssistantCreateRequest{
		Model: "gpt-4o", Name: "Math", Instructions: "tutor",
	})
	if err != nil || asst.ID != "asst_1" {
		t.Fatalf("%v %+v", err, asst)
	}
	th, err := p.CreateThread(context.Background(), &types.ThreadCreateRequest{
		Messages: []types.ThreadMessage{{Role: "user", Content: "hi"}},
	})
	if err != nil || th.ID != "thread_1" {
		t.Fatalf("%v %+v", err, th)
	}
	msg, err := p.AddMessage(context.Background(), "thread_1", types.ThreadMessage{Role: "user", Content: "2+2?"})
	if err != nil || msg.ID != "msg_1" {
		t.Fatalf("%v %+v", err, msg)
	}
	run, err := p.CreateRun(context.Background(), "thread_1", &types.RunCreateRequest{AssistantID: "asst_1"})
	if err != nil || run.ID != "run_1" {
		t.Fatalf("%v %+v", err, run)
	}
	got, err := p.GetRun(context.Background(), "thread_1", "run_1")
	if err != nil || got.Status != "completed" {
		t.Fatalf("%v %+v", err, got)
	}
}
