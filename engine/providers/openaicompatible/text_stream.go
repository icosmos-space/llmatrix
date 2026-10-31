package openaicompatible

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// TextCompletionStream implements provider.TextCompletionStreamProvider.
func (p *Provider) TextCompletionStream(ctx context.Context, req *types.TextCompletionRequest) (stream.Stream, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" || req.Prompt == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model and prompt are required", nil)
	}

	body := map[string]any{
		"model":  req.Model,
		"prompt": req.Prompt,
		"stream": true,
	}
	if req.MaxTokens != nil {
		body["max_tokens"] = *req.MaxTokens
	}
	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}
	if req.Stop != nil {
		body["stop"] = req.Stop
	}
	if req.User != "" {
		body["user"] = req.User
	}

	httpReq, err := p.newRequest(ctx, "POST", "/completions", body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")
	resp, err := p.httpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, "http request failed", err)
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, mapHTTPError(p.Info().Name, resp.StatusCode, data)
	}
	return &textSSEStream{
		provider: p.Info().Name,
		body:     resp.Body,
		scanner:  bufio.NewScanner(resp.Body),
	}, nil
}

type textSSEStream struct {
	provider string
	body     io.ReadCloser
	scanner  *bufio.Scanner
	chunk    stream.Chunk
	err      error
	done     bool
}

func (s *textSSEStream) Next(ctx context.Context) bool {
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
				s.err = errors.Wrap(errors.KindProvider, s.provider, "read stream", err)
			}
			s.done = true
			return false
		}
		line := s.scanner.Text()
		if line == "" || strings.HasPrefix(line, ":") || !strings.HasPrefix(line, "data:") {
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
				Text         string  `json:"text"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
			Error *wireError `json:"error"`
		}
		if err := json.Unmarshal([]byte(payload), &raw); err != nil {
			continue
		}
		if raw.Error != nil && raw.Error.Message != "" {
			s.err = errors.Wrap(errors.KindProvider, s.provider, raw.Error.Message, nil)
			s.done = true
			return false
		}
		ch := stream.Chunk{}
		if len(raw.Choices) > 0 {
			ch.Delta = raw.Choices[0].Text
			if raw.Choices[0].FinishReason != nil {
				ch.FinishReason = *raw.Choices[0].FinishReason
			}
		}
		s.chunk = ch
		return true
	}
}

func (s *textSSEStream) Chunk() stream.Chunk { return s.chunk }
func (s *textSSEStream) Err() error          { return s.err }
func (s *textSSEStream) Close() error {
	s.done = true
	if s.body != nil {
		return s.body.Close()
	}
	return nil
}
