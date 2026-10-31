package anthropic

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

// ChatCompletionStream implements provider.ChatStreamProvider.
func (p *Provider) ChatCompletionStream(ctx context.Context, req *types.ChatRequest) (stream.Stream, error) {
	body, err := buildRequest(req, true)
	if err != nil {
		return nil, err
	}
	httpReq, err := p.newRequest(ctx, "POST", "/v1/messages", body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := p.httpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, "anthropic", "http request failed", err)
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, mapHTTPError(resp.StatusCode, data)
	}

	return &sseStream{
		body:    resp.Body,
		scanner: bufio.NewScanner(resp.Body),
	}, nil
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
				s.err = errors.Wrap(errors.KindProvider, "anthropic", "read stream", err)
			}
			s.done = true
			return false
		}
		line := s.scanner.Text()
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if strings.HasPrefix(line, "event:") {
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
		var ev wireStreamEvent
		if err := json.Unmarshal([]byte(payload), &ev); err != nil {
			continue
		}
		if ev.Error != nil && ev.Error.Message != "" {
			s.err = errors.Wrap(errors.KindProvider, "anthropic", ev.Error.Message, nil)
			s.done = true
			return false
		}
		switch ev.Type {
		case "content_block_delta":
			if ev.Delta != nil && ev.Delta.Text != "" {
				s.chunk = stream.Chunk{Delta: ev.Delta.Text}
				return true
			}
		case "message_delta":
			if ev.Delta != nil && ev.Delta.StopReason != "" {
				s.chunk = stream.Chunk{FinishReason: mapStopReason(ev.Delta.StopReason)}
				return true
			}
		case "message_stop":
			s.chunk = stream.Chunk{Done: true}
			s.done = true
			return true
		}
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
