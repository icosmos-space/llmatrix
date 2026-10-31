package sdk

import "testing"

func TestSupportsHelpers(t *testing.T) {
	if !SupportsFunctionCalling("gpt-4o-mini") {
		t.Fatal("tools")
	}
	if !SupportsVision("gpt-4o") {
		t.Fatal("vision")
	}
	params := GetSupportedOpenAIParams("ollama/llama3")
	found := false
	for _, p := range params {
		if p == "temperature" {
			found = true
		}
		if p == "tools" {
			t.Fatal("ollama default should not list tools")
		}
	}
	if !found {
		t.Fatalf("%v", params)
	}
	info := GetModelInfo("claude-sonnet-4-5")
	if info.Provider != "anthropic" || !info.SupportsPromptCaching {
		t.Fatalf("%+v", info)
	}
}
