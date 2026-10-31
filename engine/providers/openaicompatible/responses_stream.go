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

// ResponsesStream implements provider.ResponsesStreamProvider via POST /responses?stream.
func (p *Provider) ResponsesStream(ctx context.Context, req *types.ResponsesRequest) (stream.Stream, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	if req.Input == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "input is required", nil)
	}

	body := responsesBody(req, true)
	httpReq, err := p.newRequest(ctx, "POST", "/responses", body)
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

	return &responsesSSEStream{
		provider: p.Info().Name,
		body:     resp.Body,
		scanner:  bufio.NewScanner(resp.Body),
	}, nil
}

type responsesSSEStream struct {
	provider string
	body     io.ReadCloser
	scanner  *bufio.Scanner
	chunk    stream.Chunk
	err      error
	done     bool
}

func (s *responsesSSEStream) Next(ctx context.Context) bool {
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
		var ev struct {
			Type  string `json:"type"`
			Delta string `json:"delta"`
			Text  string `json:"text"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(payload), &ev); err != nil {
			continue
		}
		if ev.Error != nil && ev.Error.Message != "" {
			s.err = errors.Wrap(errors.KindProvider, s.provider, ev.Error.Message, nil)
			s.done = true
			return false
		}
		switch {
		case strings.Contains(ev.Type, "output_text.delta"), ev.Type == "response.output_text.delta":
			delta := ev.Delta
			if delta == "" {
				delta = ev.Text
			}
			s.chunk = stream.Chunk{Delta: delta}
			return true
		case strings.HasSuffix(ev.Type, ".completed"), ev.Type == "response.completed", ev.Type == "response.done":
			s.chunk = stream.Chunk{Done: true, FinishReason: "stop"}
			s.done = true
			return true
		default:
			// Ignore lifecycle events (created, in_progress, content_part.*, etc.).
			continue
		}
	}
}

func (s *responsesSSEStream) Chunk() stream.Chunk { return s.chunk }
func (s *responsesSSEStream) Err() error          { return s.err }
func (s *responsesSSEStream) Close() error {
	s.done = true
	if s.body != nil {
		return s.body.Close()
	}
	return nil
}
