package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
)

// Client talks to a deployed LLMatrix Proxy over HTTP.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a Proxy client.
func NewClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 120 * time.Second}
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// ChatCompletion sends a chat request to the proxy.
func (c *Client) ChatCompletion(ctx context.Context, req *chat.Request) (*chat.Response, error) {
	var out chat.Response
	if err := c.doJSON(ctx, http.MethodPost, "/v1/chat/completions", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Embeddings sends an embedding request to the proxy.
func (c *Client) Embeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
	var out embedding.Response
	if err := c.doJSON(ctx, http.MethodPost, "/v1/embeddings", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, body, dest any) error {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("llmatrix/sdk/proxy: marshal: %w", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, rdr)
	if err != nil {
		return fmt.Errorf("llmatrix/sdk/proxy: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("llmatrix/sdk/proxy: http: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("llmatrix/sdk/proxy: read: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("llmatrix/sdk/proxy: HTTP %d: %s", resp.StatusCode, truncate(string(data), 500))
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("llmatrix/sdk/proxy: decode: %w", err)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
