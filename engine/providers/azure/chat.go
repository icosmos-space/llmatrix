package azure

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

type wireToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type wireMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
	ToolCalls  []wireToolCall `json:"tool_calls,omitempty"`
}

type wireChatRequest struct {
	Messages    []wireMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float64      `json:"top_p,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Tools       any           `json:"tools,omitempty"`
	ToolChoice  any           `json:"tool_choice,omitempty"`
	User        string        `json:"user,omitempty"`
	Seed        *int          `json:"seed,omitempty"`
	Stop        []string      `json:"stop,omitempty"`
}

type wireChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      wireMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage types.Usage `json:"usage"`
}

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	if req == nil || req.Model == "" || len(req.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "azure", "model and messages are required", nil)
	}
	body := wireChatRequest{
		Messages:    toWireMessages(req.Messages),
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
		User:        req.User,
		Seed:        req.Seed,
		Stop:        req.Stop,
		ToolChoice:  req.ToolChoice,
	}
	if len(req.Tools) > 0 {
		body.Tools = req.Tools
	}
	var raw wireChatResponse
	if err := p.doJSON(ctx, "POST", req.Model, "/chat/completions", body, &raw); err != nil {
		return nil, err
	}
	out := &types.ChatResponse{
		ID:       raw.ID,
		Model:    firstNonEmpty(raw.Model, req.Model),
		Provider: "azure",
		Usage:    raw.Usage,
		Choices:  make([]types.ChatChoice, 0, len(raw.Choices)),
	}
	for _, c := range raw.Choices {
		msg := types.Message{Role: types.Role(c.Message.Role), Content: c.Message.Content, ToolCallID: c.Message.ToolCallID}
		for _, tc := range c.Message.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, types.ToolCall{
				ID: tc.ID, Type: tc.Type, Name: tc.Function.Name, Arguments: tc.Function.Arguments,
			})
		}
		out.Choices = append(out.Choices, types.ChatChoice{Index: c.Index, Message: msg, FinishReason: c.FinishReason})
	}
	return out, nil
}

// ChatCompletionStream implements provider.ChatStreamProvider.
func (p *Provider) ChatCompletionStream(ctx context.Context, req *types.ChatRequest) (stream.Stream, error) {
	if req == nil || req.Model == "" || len(req.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "azure", "model and messages are required", nil)
	}
	body := wireChatRequest{
		Messages:    toWireMessages(req.Messages),
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      true,
	}
	rawURL, err := p.deploymentURL(req.Model, "/chat/completions")
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, "azure", "marshal request", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", rawURL, bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "azure", "create request", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if key := p.resolvedAPIKey(); key != "" {
		httpReq.Header.Set("api-key", key)
	}
	resp, err := p.httpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, "azure", "http request failed", err)
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, mapHTTPError(resp.StatusCode, data)
	}
	return &sseStream{body: resp.Body, scanner: bufio.NewScanner(resp.Body)}, nil
}

// Embeddings implements provider.EmbeddingProvider.
func (p *Provider) Embeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	if req == nil || req.Model == "" || len(req.Input) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "azure", "model and input are required", nil)
	}
	body := map[string]any{"input": req.Input}
	var raw struct {
		Data []struct {
			Index     int       `json:"index"`
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
		Usage types.Usage `json:"usage"`
		Model string      `json:"model"`
	}
	if err := p.doJSON(ctx, "POST", req.Model, "/embeddings", body, &raw); err != nil {
		return nil, err
	}
	out := &types.EmbeddingResponse{
		Model: firstNonEmpty(raw.Model, req.Model), Provider: "azure", Usage: raw.Usage,
		Data: make([]types.EmbeddingItem, 0, len(raw.Data)),
	}
	for _, d := range raw.Data {
		out.Data = append(out.Data, types.EmbeddingItem{Index: d.Index, Embedding: d.Embedding})
	}
	return out, nil
}

func toWireMessages(msgs []types.Message) []wireMessage {
	out := make([]wireMessage, 0, len(msgs))
	for _, m := range msgs {
		wm := wireMessage{Role: string(m.Role), Content: m.Content, ToolCallID: m.ToolCallID}
		for _, tc := range m.ToolCalls {
			wtc := wireToolCall{ID: tc.ID, Type: firstNonEmpty(tc.Type, "function")}
			wtc.Function.Name = tc.Name
			wtc.Function.Arguments = tc.Arguments
			wm.ToolCalls = append(wm.ToolCalls, wtc)
		}
		out = append(out, wm)
	}
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

type sseStream struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
	chunk   stream.Chunk
	err     error
	done    bool
}

func (s *sseStream) Next(ctx context.Context) bool {
	if s.done || s.err != nil {
		return false
	}
	for {
		select {
		case <-ctx.Done():
			s.err = ctx.Err()
			s.done = true
			return false
		default:
		}
		if !s.scanner.Scan() {
			if err := s.scanner.Err(); err != nil {
				s.err = err
			}
			s.done = true
			return false
		}
		line := s.scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			s.chunk = stream.Chunk{Done: true}
			s.done = true
			return true
		}
		var raw struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(payload), &raw); err != nil {
			continue
		}
		ch := stream.Chunk{}
		if len(raw.Choices) > 0 {
			ch.Delta = raw.Choices[0].Delta.Content
			if raw.Choices[0].FinishReason != nil {
				ch.FinishReason = *raw.Choices[0].FinishReason
			}
		}
		s.chunk = ch
		return true
	}
}

func (s *sseStream) Chunk() stream.Chunk { return s.chunk }
func (s *sseStream) Err() error          { return s.err }
func (s *sseStream) Close() error {
	s.done = true
	if s.body != nil {
		return s.body.Close()
	}
	return nil
}
