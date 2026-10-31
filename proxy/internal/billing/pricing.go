package billing

// Price describes per-1K-token rates.
type Price struct {
	InputPer1K  float64
	OutputPer1K float64
}

// PricingTable stores model prices keyed by provider/model.
type PricingTable struct {
	prices map[string]Price
}

// NewPricingTable creates an empty table.
func NewPricingTable() *PricingTable {
	return &PricingTable{prices: make(map[string]Price)}
}

func key(provider, model string) string { return provider + "/" + model }

// Set registers a price.
func (t *PricingTable) Set(provider, model string, p Price) {
	t.prices[key(provider, model)] = p
}

// Get looks up a price.
func (t *PricingTable) Get(provider, model string) (Price, bool) {
	p, ok := t.prices[key(provider, model)]
	return p, ok
}
