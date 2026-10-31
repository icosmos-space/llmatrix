package router

import "testing"

func TestFallbackCandidates(t *testing.T) {
	s := FallbackStrategy{
		Primary:   ParseModelRef("openai/gpt-4o"),
		Fallbacks: ParseModelRefs("deepseek/deepseek-chat", "ollama/llama3"),
	}
	cands := s.Candidates()
	if len(cands) != 3 {
		t.Fatalf("len=%d", len(cands))
	}
	if cands[0].Provider != "openai" || cands[1].Model != "deepseek-chat" {
		t.Fatalf("%+v", cands)
	}
}
