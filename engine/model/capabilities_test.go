package model

import "testing"

func TestSupportsHelpers(t *testing.T) {
	if !SupportsFunctionCalling("gpt-4o") {
		t.Fatal("gpt-4o tools")
	}
	if !SupportsVision("openai/gpt-4o-mini") {
		t.Fatal("vision")
	}
	if SupportsVision("openai/gpt-3.5-turbo-instruct") {
		t.Fatal("instruct should not be vision")
	}
	params := GetSupportedOpenAIParams("anthropic/claude-sonnet-4-5")
	if len(params) < 5 {
		t.Fatalf("%v", params)
	}
	info := GetModelInfo("deepseek-chat", "")
	if info.Provider != "deepseek" || !info.SupportsFunctionCalling {
		t.Fatalf("%+v", info)
	}
}
