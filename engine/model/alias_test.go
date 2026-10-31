package model

import "testing"

func TestResolveAlias(t *testing.T) {
	a := DefaultAliases()
	ref := ResolveAlias(a, "gpt-4o-mini", "")
	if ref.Provider != "openai" || ref.Model != "gpt-4o-mini" {
		t.Fatalf("%+v", ref)
	}
	ref = ResolveAlias(a, "anthropic/claude-sonnet-4-5", "")
	if ref.Provider != "anthropic" {
		t.Fatalf("%+v", ref)
	}
}
