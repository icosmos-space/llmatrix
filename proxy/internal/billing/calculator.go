package billing

import "github.com/icosmos-space/llmatrix/engine/types"

// Calculator computes cost from usage and pricing tables.
type Calculator struct {
	Pricing *PricingTable
}

// NewCalculator creates a Calculator.
func NewCalculator(pricing *PricingTable) *Calculator {
	return &Calculator{Pricing: pricing}
}

// Cost returns the USD cost for the given usage.
func (c *Calculator) Cost(provider, model string, u types.Usage) float64 {
	if c.Pricing == nil {
		return 0
	}
	p, ok := c.Pricing.Get(provider, model)
	if !ok {
		return 0
	}
	in := float64(u.PromptTokens) / 1000.0 * p.InputPer1K
	out := float64(u.CompletionTokens) / 1000.0 * p.OutputPer1K
	return in + out
}
