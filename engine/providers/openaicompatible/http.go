package openaicompatible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/icosmos-space/llmatrix/engine/errors"
)

func (p *Provider) httpClient() *http.Client {
	if p.HTTPClient != nil {
		return p.HTTPClient
	}
	return &http.Client{Timeout: 120 * time.Second}
}

func (p *Provider) endpoint(path string) (string, error) {
	base := strings.TrimRight(p.resolvedBaseURL(), "/")
	if base == "" {
		return "", errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "base URL is required", nil)
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path, nil
}

func (p *Provider) resolvedBaseURL() string {
	return ResolveBaseURL(p.Info().Name, p.BaseURL)
}

func (p *Provider) resolvedAPIKey() string {
	return ResolveAPIKey(p.Info().Name, p.APIKey)
}

func (p *Provider) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	url, err := p.endpoint(path)
	if err != nil {
		return nil, err
	}
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "marshal request", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	p.applyAuthHeaders(req)
	return req, nil
}

func (p *Provider) applyAuthHeaders(req *http.Request) {
	if key := p.resolvedAPIKey(); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	if p.Organization != "" {
		req.Header.Set("OpenAI-Organization", p.Organization)
	}
	for k, v := range p.ExtraHeaders {
		if k == "" || v == "" {
			continue
		}
		req.Header.Set(k, v)
	}
}

func (p *Provider) doJSON(ctx context.Context, method, path string, body, dest any) error {
	req, err := p.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	resp, err := p.httpClient().Do(req)
	if err != nil {
		return errors.Wrap(errors.KindProvider, p.Info().Name, "http request failed", err)
	}
	defer resp.Body.Close()
	return p.decodeJSONResponse(resp, dest)
}

func (p *Provider) decodeJSONResponse(resp *http.Response, dest any) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.KindProvider, p.Info().Name, "read response", err)
	}
	if resp.StatusCode >= 400 {
		return mapHTTPError(p.Info().Name, resp.StatusCode, data)
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return errors.Wrap(errors.KindProvider, p.Info().Name, fmt.Sprintf("decode response: %s", truncate(string(data), 200)), err)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
