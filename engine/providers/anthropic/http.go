package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
)

func (p *Provider) endpoint(path string) string {
	base := strings.TrimRight(p.resolvedBaseURL(), "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

func (p *Provider) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(errors.KindInvalidRequest, "anthropic", "marshal request", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.endpoint(path), rdr)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "anthropic", "create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("anthropic-version", p.resolvedVersion())
	if key := p.resolvedAPIKey(); key != "" {
		req.Header.Set("x-api-key", key)
	}
	return req, nil
}

func (p *Provider) doJSON(ctx context.Context, method, path string, body, dest any) error {
	req, err := p.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	resp, err := p.httpClient().Do(req)
	if err != nil {
		return errors.Wrap(errors.KindProvider, "anthropic", "http request failed", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.KindProvider, "anthropic", "read response", err)
	}
	if resp.StatusCode >= 400 {
		return mapHTTPError(resp.StatusCode, data)
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return errors.Wrap(errors.KindProvider, "anthropic", fmt.Sprintf("decode response: %s", truncate(string(data), 200)), err)
	}
	return nil
}

func mapHTTPError(status int, body []byte) error {
	msg := string(body)
	var payload struct {
		Error *wireError `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Error != nil && payload.Error.Message != "" {
		msg = payload.Error.Message
	}
	kind := errors.KindProvider
	switch status {
	case http.StatusBadRequest:
		kind = errors.KindInvalidRequest
	case http.StatusUnauthorized, http.StatusForbidden:
		kind = errors.KindUnauthorized
	case http.StatusNotFound:
		kind = errors.KindNotFound
	case http.StatusTooManyRequests:
		kind = errors.KindRateLimited
	case http.StatusRequestTimeout, 529:
		kind = errors.KindTimeout
	}
	return errors.WrapStatus(kind, "anthropic", fmt.Sprintf("HTTP %d: %s", status, truncate(msg, 500)), status, nil)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
