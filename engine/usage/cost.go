package usage

import "github.com/icosmos-space/llmatrix/engine/types"

// Price is USD cost per 1K tokens.
type Price struct {
	InputPer1K  float64
	OutputPer1K float64
}

// DefaultPrices is a small built-in table (approximate; override in production).
var DefaultPrices = map[string]Price{
	"openai/gpt-4o-mini":            {InputPer1K: 0.00015, OutputPer1K: 0.0006},
	"openai/gpt-4o":                 {InputPer1K: 0.0025, OutputPer1K: 0.01},
	"openai/text-embedding-3-small": {InputPer1K: 0.00002, OutputPer1K: 0},
	"anthropic/claude-sonnet-4-5":   {InputPer1K: 0.003, OutputPer1K: 0.015},
	"deepseek/deepseek-chat":        {InputPer1K: 0.00014, OutputPer1K: 0.00028},
	"gemini/gemini-2.0-flash":       {InputPer1K: 0.0001, OutputPer1K: 0.0004},
}

// CostUSD estimates USD cost for usage of a provider/model key.
func CostUSD(providerModel string, u types.Usage, table map[string]Price) float64 {
	if table == nil {
		table = DefaultPrices
	}
	p, ok := table[providerModel]
	if !ok {
		return 0
	}
	return float64(u.PromptTokens)/1000*p.InputPer1K + float64(u.CompletionTokens)/1000*p.OutputPer1K
}
