package proxy

import (
	"context"
	"net/http"

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
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// ChatCompletion sends a chat request to the proxy.
func (c *Client) ChatCompletion(ctx context.Context, req *chat.Request) (*chat.Response, error) {
	_ = ctx
	_ = req
	return nil, errNotImplemented
}

// Embeddings sends an embedding request to the proxy.
func (c *Client) Embeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
	_ = ctx
	_ = req
	return nil, errNotImplemented
}
