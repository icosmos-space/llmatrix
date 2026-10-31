package sdk

import (
	"github.com/icosmos-space/llmatrix/engine/cache"
)

func (c *Client) cacheEnabled(flag *bool) bool {
	return flag != nil && *flag && c.cfg.Cache != nil
}

func (c *Client) cacheGet(method, model string, payload any) (any, bool) {
	if c.cfg.Cache == nil {
		return nil, false
	}
	return c.cfg.Cache.Get(cache.Key(method, model, payload))
}

func (c *Client) cacheSet(method, model string, payload, value any) {
	if c.cfg.Cache == nil {
		return
	}
	c.cfg.Cache.Set(cache.Key(method, model, payload), value, c.cfg.CacheTTL)
}

func (c *Client) cacheModelKey(requestModel string) string {
	refs := c.candidates(requestModel)
	if len(refs) == 0 {
		return requestModel
	}
	return refs[0].Provider + "/" + refs[0].Model
}
