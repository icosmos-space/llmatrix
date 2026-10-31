package model

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in       string
		provider string
		model    string
	}{
		{"openai/gpt-4o", "openai", "gpt-4o"},
		{"deepseek/deepseek-chat", "deepseek", "deepseek-chat"},
		{"gpt-4o", "", "gpt-4o"},
		{"  ollama/llama3  ", "ollama", "llama3"},
	}
	for _, tc := range cases {
		got := Parse(tc.in)
		if got.Provider != tc.provider || got.Model != tc.model {
			t.Fatalf("Parse(%q)=%+v want provider=%q model=%q", tc.in, got, tc.provider, tc.model)
		}
	}
}

func TestResolveDefaultsProvider(t *testing.T) {
	ref := Resolve("gpt-4o-mini", "")
	if ref.Provider != "openai" || ref.Model != "gpt-4o-mini" {
		t.Fatalf("unexpected %+v", ref)
	}
	ref = Resolve("deepseek-chat", "deepseek")
	if ref.Provider != "deepseek" || ref.Model != "deepseek-chat" {
		t.Fatalf("unexpected %+v", ref)
	}
}
