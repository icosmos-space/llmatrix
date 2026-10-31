package integration_test

import (
	"testing"

	"github.com/icosmos-space/llmatrix/engine/providers/openai"
	"github.com/icosmos-space/llmatrix/engine/registry"
)

func TestRegistryRegistersOpenAI(t *testing.T) {
	reg := registry.New()
	if err := reg.Register(openai.New("")); err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Get("openai"); !ok {
		t.Fatal("expected openai provider")
	}
}
