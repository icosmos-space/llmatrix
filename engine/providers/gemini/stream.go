package gemini

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
// Uses streamGenerateContent with alt=sse.
func (p *Provider) ChatCompletionStream(ctx context.Context, req *types.ChatRequest) (stream.Stream, error) {
	body, model, err := buildRequest(req)
	if err != nil {
		return nil, err
	}
	rawURL, err := p.modelURL(model, "streamGenerateContent")
	if err != nil {
		return nil, err
	}
	if strings.Contains(rawURL, "?") {
		rawURL += "&alt=sse"
	} else {
		rawURL += "?alt=sse"
	}

	httpReq, err := p.newRequest(ctx, "POST", rawURL, body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := p.httpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, "gemini", "http request failed", err)
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
				s.err = errors.Wrap(errors.KindProvider, "gemini", "read stream", err)
			} else {
				s.chunk = stream.Chunk{Done: true}
				s.done = true
				return true
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
		var raw wireResponse
		if err := json.Unmarshal([]byte(payload), &raw); err != nil {
			continue
		}
		if raw.Error != nil && raw.Error.Message != "" {
			s.err = errors.Wrap(errors.KindProvider, "gemini", raw.Error.Message, nil)
			s.done = true
			return false
		}
		if len(raw.Candidates) == 0 {
			continue
		}
		cand := raw.Candidates[0]
		text := partsText(cand.Content)
		ch := stream.Chunk{Delta: text}
		if cand.FinishReason != "" {
			ch.FinishReason = mapFinishReason(cand.FinishReason)
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
