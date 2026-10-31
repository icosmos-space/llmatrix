package usage

import "github.com/icosmos-space/llmatrix/engine/types"

// Price is USD cost per 1K tokens.
type Price struct {
	InputPer1K  float64
	OutputPer1K float64
}

// DefaultPrices is a small built-in table (approximate; override in production).
var DefaultPrices = map[string]Price{
	"openai/gpt-4o-mini":                 {InputPer1K: 0.00015, OutputPer1K: 0.0006},
	"openai/gpt-4o":                      {InputPer1K: 0.0025, OutputPer1K: 0.01},
	"openai/gpt-4.1":                     {InputPer1K: 0.002, OutputPer1K: 0.008},
	"openai/gpt-4.1-mini":                {InputPer1K: 0.0004, OutputPer1K: 0.0016},
	"openai/o1-mini":                     {InputPer1K: 0.0011, OutputPer1K: 0.0044},
	"openai/o3-mini":                     {InputPer1K: 0.0011, OutputPer1K: 0.0044},
	"openai/text-embedding-3-small":      {InputPer1K: 0.00002, OutputPer1K: 0},
	"openai/text-embedding-3-large":      {InputPer1K: 0.00013, OutputPer1K: 0},
	"anthropic/claude-sonnet-4-5":        {InputPer1K: 0.003, OutputPer1K: 0.015},
	"anthropic/claude-opus-4":            {InputPer1K: 0.015, OutputPer1K: 0.075},
	"anthropic/claude-3-5-sonnet-latest": {InputPer1K: 0.003, OutputPer1K: 0.015},
	"deepseek/deepseek-chat":             {InputPer1K: 0.00014, OutputPer1K: 0.00028},
	"deepseek/deepseek-reasoner":         {InputPer1K: 0.00055, OutputPer1K: 0.00219},
	"gemini/gemini-2.0-flash":            {InputPer1K: 0.0001, OutputPer1K: 0.0004},
	"gemini/gemini-1.5-pro":              {InputPer1K: 0.00125, OutputPer1K: 0.005},
	"groq/llama-3.3-70b-versatile":       {InputPer1K: 0.00059, OutputPer1K: 0.00079},
	"mistral/mistral-large-latest":       {InputPer1K: 0.002, OutputPer1K: 0.006},
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
