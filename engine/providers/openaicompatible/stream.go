package openaicompatible

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatCompletionStream implements provider.ChatStreamProvider.
func (p *Provider) ChatCompletionStream(ctx context.Context, req *types.ChatRequest) (stream.Stream, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	if len(req.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "messages are required", nil)
	}

	body := toWireChatRequest(req, true)

	httpReq, err := p.newRequest(ctx, "POST", "/chat/completions", body)
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

	return &sseStream{
		provider: p.Info().Name,
		body:     resp.Body,
		scanner:  bufio.NewScanner(resp.Body),
	}, nil
}

type sseStream struct {
	provider string
	body     io.ReadCloser
	scanner  *bufio.Scanner
	chunk    stream.Chunk
	err      error
	done     bool
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
				s.err = errors.Wrap(errors.KindProvider, s.provider, "read stream", err)
			}
			s.done = true
			return false
		}
		line := s.scanner.Text()
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			s.chunk = stream.Chunk{Done: true}
			s.done = true
			return true
		}
		var raw wireStreamChunk
		if err := json.Unmarshal([]byte(payload), &raw); err != nil {
			// Some providers send concatenated JSON; try to recover next event.
			if bytes.Contains([]byte(payload), []byte("{")) {
				s.err = errors.Wrap(errors.KindProvider, s.provider, "decode stream chunk", err)
				s.done = true
				return false
			}
			continue
		}
		if raw.Error != nil && raw.Error.Message != "" {
			s.err = errors.Wrap(errors.KindProvider, s.provider, raw.Error.Message, nil)
			s.done = true
			return false
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

func (s *sseStream) Err() error { return s.err }

func (s *sseStream) Close() error {
	s.done = true
	if s.body != nil {
		return s.body.Close()
	}
	return nil
}
