package sdk

import (
	"testing"

	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func TestTokenCounterAndCost(t *testing.T) {
	if TokenCounter("hello world") < 1 {
		t.Fatal("token counter")
	}
	n := TokenCounterMessages([]chat.Message{{Role: "user", Content: "hi there"}})
	if n < 1 {
		t.Fatalf("messages=%d", n)
	}
	cost := CompletionCost("gpt-4o-mini", 1000, 1000)
	if cost <= 0 {
		t.Fatalf("cost=%v", cost)
	}
}
