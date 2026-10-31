package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/icosmos-space/llmatrix/engine/callback"
	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/registry"
	"github.com/icosmos-space/llmatrix/engine/retry"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Client is the unified entry point for LLMatrix SDK calls
// (LiteLLM-style: one API, multiple providers).
type Client struct {
	cfg Config
}

// NewClient constructs a Client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{cfg: defaultConfig()}
	for _, opt := range opts {
		opt(&c.cfg)
	}
	return c
}

func (c *Client) candidates(requestModel string) []model.Ref {
	rawModel := requestModel
	if rawModel == "" {
		rawModel = c.cfg.Model
	}
	if rawModel == "" {
		return nil
	}
	primary := model.ResolveAlias(c.cfg.Aliases, rawModel, c.cfg.Provider)
	out := []model.Ref{primary}
	for _, fb := range c.cfg.Fallbacks {
		ref := model.ResolveAlias(c.cfg.Aliases, fb, "")
		if ref.Provider == primary.Provider && ref.Model == primary.Model {
			continue
		}
		out = append(out, ref)
	}
	return out
}

func (c *Client) newProvider(ref model.Ref, shareClientCreds bool) (provider.Provider, error) {
	cfg := registry.ProviderConfig{
		Name:       ref.Provider,
		HTTPClient: c.cfg.httpClient(),
		APIVersion: c.cfg.APIVersion,
	}
	if shareClientCreds {
		cfg.APIKey = c.cfg.APIKey
		cfg.BaseURL = c.cfg.BaseURL
	}
	return registry.NewProvider(cfg)
}

type candidateFn func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error)

func (c *Client) withCandidates(ctx context.Context, method, requestModel string, fn candidateFn) (any, error) {
	candidates := c.candidates(requestModel)
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: model is required (set request.Model or WithModel)", ErrInvalidConfig)
	}
	primaryProvider := candidates[0].Provider
	var last error
	for i, ref := range candidates {
		start := time.Now()
		p, err := c.newProvider(ref, ref.Provider == primaryProvider)
		if err != nil {
			last = MapError(err)
			continue
		}
		var usage types.Usage
		var out any
		err = retry.Do(ctx, c.cfg.Retry, func(ctx context.Context) error {
			u, v, callErr := fn(ctx, ref, p)
			if callErr != nil {
				return callErr
			}
			usage = u
			out = v
			return nil
		})
		ev := callback.Event{
			Method:   method,
			Provider: ref.Provider,
			Model:    ref.Model,
			Start:    start,
			End:      time.Now(),
			Usage:    usage,
		}
		if err == nil {
			if c.cfg.Callbacks != nil {
				c.cfg.Callbacks.EmitSuccess(ctx, ev)
			}
			return out, nil
		}
		mapped := MapError(err)
		ev.Err = mapped
		if c.cfg.Callbacks != nil {
			c.cfg.Callbacks.EmitFailure(ctx, ev)
		}
		last = mapped
		if i < len(candidates)-1 && engineerrors.IsFallbackable(err) {
			continue
		}
		return nil, mapped
	}
	return nil, last
}
