package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/provider"
)

const defaultAPIVersion = "2024-08-01-preview"

// Provider implements Azure OpenAI (deployment-based OpenAI-compatible API).
// Model is the deployment name. BaseURL is https://{resource}.openai.azure.com
type Provider struct {
	APIKey     string
	BaseURL    string
	APIVersion string
	HTTPClient *http.Client
}

// New creates an Azure OpenAI provider.
func New(apiKey string, opts ...Option) *Provider {
	p := &Provider{APIKey: apiKey, APIVersion: defaultAPIVersion}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures Azure.
type Option func(*Provider)

// WithBaseURL sets the Azure resource endpoint (without /openai/deployments/...).
func WithBaseURL(u string) Option {
	return func(p *Provider) {
		if u != "" {
			p.BaseURL = strings.TrimRight(u, "/")
		}
	}
}

// WithAPIVersion sets the api-version query parameter.
func WithAPIVersion(v string) Option {
	return func(p *Provider) {
		if v != "" {
			p.APIVersion = v
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	return provider.Info{
		Name:        "azure",
		DisplayName: "Azure OpenAI",
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true, Embeddings: true,
			Images: true, Transcription: true,
		},
	}
}

func (p *Provider) httpClient() *http.Client {
	if p.HTTPClient != nil {
		return p.HTTPClient
	}
	return &http.Client{Timeout: 120 * time.Second}
}

func (p *Provider) resolvedAPIKey() string {
	if p.APIKey != "" {
		return p.APIKey
	}
	if v := os.Getenv("AZURE_API_KEY"); v != "" {
		return v
	}
	return os.Getenv("AZURE_OPENAI_API_KEY")
}

func (p *Provider) resolvedBaseURL() string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	if v := os.Getenv("AZURE_API_BASE"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return strings.TrimRight(os.Getenv("AZURE_OPENAI_ENDPOINT"), "/")
}

func (p *Provider) resolvedAPIVersion() string {
	if p.APIVersion != "" {
		return p.APIVersion
	}
	if v := os.Getenv("AZURE_API_VERSION"); v != "" {
		return v
	}
	return defaultAPIVersion
}

func (p *Provider) deploymentURL(deployment, path string) (string, error) {
	base := p.resolvedBaseURL()
	if base == "" {
		return "", errors.Wrap(errors.KindInvalidRequest, "azure", "base URL is required (AZURE_API_BASE)", nil)
	}
	if deployment == "" {
		return "", errors.Wrap(errors.KindInvalidRequest, "azure", "deployment (model) is required", nil)
	}
	u, err := url.Parse(fmt.Sprintf("%s/openai/deployments/%s%s", base, url.PathEscape(deployment), path))
	if err != nil {
		return "", errors.Wrap(errors.KindInternal, "azure", "parse url", err)
	}
	q := u.Query()
	q.Set("api-version", p.resolvedAPIVersion())
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (p *Provider) doJSON(ctx context.Context, method, deployment, path string, body, dest any) error {
	rawURL, err := p.deploymentURL(deployment, path)
	if err != nil {
		return err
	}
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(errors.KindInvalidRequest, "azure", "marshal request", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, rdr)
	if err != nil {
		return errors.Wrap(errors.KindInternal, "azure", "create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if key := p.resolvedAPIKey(); key != "" {
		req.Header.Set("api-key", key)
	}
	resp, err := p.httpClient().Do(req)
	if err != nil {
		return errors.Wrap(errors.KindProvider, "azure", "http request failed", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.KindProvider, "azure", "read response", err)
	}
	if resp.StatusCode >= 400 {
		return mapHTTPError(resp.StatusCode, data)
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return errors.Wrap(errors.KindProvider, "azure", "decode response", err)
	}
	return nil
}

func mapHTTPError(status int, body []byte) error {
	msg := string(body)
	var payload struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
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
	return errors.WrapStatus(kind, "azure", fmt.Sprintf("HTTP %d: %s", status, truncate(msg, 500)), status, nil)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
