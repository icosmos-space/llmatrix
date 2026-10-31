package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
)

func (p *Provider) modelURL(model, action string) (string, error) {
	base := strings.TrimRight(p.resolvedBaseURL(), "/")
	model = strings.TrimPrefix(model, "models/")
	if model == "" {
		return "", errors.Wrap(errors.KindInvalidRequest, "gemini", "model is required", nil)
	}
	u, err := url.Parse(fmt.Sprintf("%s/models/%s:%s", base, model, action))
	if err != nil {
		return "", errors.Wrap(errors.KindInternal, "gemini", "parse url", err)
	}
	q := u.Query()
	if key := p.resolvedAPIKey(); key != "" {
		q.Set("key", key)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (p *Provider) newRequest(ctx context.Context, method, rawURL string, body any) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(errors.KindInvalidRequest, "gemini", "marshal request", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, rdr)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "gemini", "create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (p *Provider) doJSON(ctx context.Context, method, rawURL string, body, dest any) error {
	req, err := p.newRequest(ctx, method, rawURL, body)
	if err != nil {
		return err
	}
	resp, err := p.httpClient().Do(req)
	if err != nil {
		return errors.Wrap(errors.KindProvider, "gemini", "http request failed", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.KindProvider, "gemini", "read response", err)
	}
	if resp.StatusCode >= 400 {
		return mapHTTPError(resp.StatusCode, data)
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return errors.Wrap(errors.KindProvider, "gemini", fmt.Sprintf("decode response: %s", truncate(string(data), 200)), err)
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
	}
	return errors.WrapStatus(kind, "gemini", fmt.Sprintf("HTTP %d: %s", status, truncate(msg, 500)), status, nil)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
