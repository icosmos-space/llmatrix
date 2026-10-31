package sdk

import (
	"context"
	"fmt"

	enginerouter "github.com/icosmos-space/llmatrix/engine/router"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
)

// Deployment is a LiteLLM model_list-style upstream entry.
type Deployment = enginerouter.Deployment

// Router provides LiteLLM Router-style load balancing across deployments.
type Router struct {
	lb     *enginerouter.SimpleLoadBalance
	client *Client
}

// RouterOption configures a Router.
type RouterOption func(*Router)

// WithRouterClientOptions applies Client options used for upstream calls.
func WithRouterClientOptions(opts ...Option) RouterOption {
	return func(r *Router) {
		r.client = NewClient(opts...)
	}
}

// NewRouter builds a Router from deployments (LiteLLM model_list).
func NewRouter(deployments []Deployment, opts ...RouterOption) *Router {
	r := &Router{
		lb:     enginerouter.NewSimpleLoadBalance(),
		client: NewClient(),
	}
	for _, d := range deployments {
		r.lb.Register(d)
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// ChatCompletion picks a deployment for the logical model name then calls chat.
func (r *Router) ChatCompletion(ctx context.Context, logicalModel string, req *chat.Request) (*chat.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	d, err := r.lb.Pick(logicalModel)
	if err != nil {
		return nil, err
	}
	cp := *req
	cp.Model = d.Provider + "/" + d.Model
	opts := []Option{}
	if d.APIKey != "" {
		opts = append(opts, WithAPIKey(d.APIKey))
	}
	if d.BaseURL != "" {
		opts = append(opts, WithBaseURL(d.BaseURL))
	}
	if len(opts) == 0 {
		return r.client.ChatCompletion(ctx, &cp)
	}
	// Merge router deployment credentials into an ephemeral client.
	cfg := r.client.cfg
	c := &Client{cfg: cfg}
	for _, opt := range opts {
		opt(&c.cfg)
	}
	return c.ChatCompletion(ctx, &cp)
}

// Embeddings picks a deployment then calls embeddings.
func (r *Router) Embeddings(ctx context.Context, logicalModel string, req *embedding.Request) (*embedding.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	d, err := r.lb.Pick(logicalModel)
	if err != nil {
		return nil, fmt.Errorf("router: %w", err)
	}
	cp := *req
	cp.Model = d.Provider + "/" + d.Model
	cfg := r.client.cfg
	c := &Client{cfg: cfg}
	if d.APIKey != "" {
		c.cfg.APIKey = d.APIKey
	}
	if d.BaseURL != "" {
		c.cfg.BaseURL = d.BaseURL
	}
	return c.Embeddings(ctx, &cp)
}
