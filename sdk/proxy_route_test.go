package sdk

import "testing"

func TestUseLiteLLMProxyCandidates(t *testing.T) {
	c := NewClient(WithUseLiteLLMProxy(true), WithModel("openai/gpt-4o-mini"))
	refs := c.candidates("")
	if len(refs) != 1 || refs[0].Provider != "litellm_proxy" || refs[0].Model != "gpt-4o-mini" {
		t.Fatalf("%+v", refs)
	}
}
