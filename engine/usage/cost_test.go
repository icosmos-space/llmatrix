package usage

import (
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestCostUSD(t *testing.T) {
	cost := CostUSD("openai/gpt-4o-mini", types.Usage{PromptTokens: 1000, CompletionTokens: 1000}, nil)
	if cost <= 0 {
		t.Fatalf("cost=%v", cost)
	}
	if CostUSD("unknown/model", types.Usage{PromptTokens: 100}, nil) != 0 {
		t.Fatal("expected 0 for unknown")
	}
}
