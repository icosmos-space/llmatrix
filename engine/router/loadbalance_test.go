package router

import "testing"

func TestSimpleLoadBalanceRoundRobin(t *testing.T) {
	lb := NewSimpleLoadBalance()
	lb.Register(Deployment{Name: "gpt", Provider: "openai", Model: "a", Weight: 1})
	lb.Register(Deployment{Name: "gpt", Provider: "openai", Model: "b", Weight: 1})

	seen := map[string]int{}
	for i := 0; i < 10; i++ {
		d, err := lb.Pick("gpt")
		if err != nil {
			t.Fatal(err)
		}
		seen[d.Model]++
	}
	if seen["a"] == 0 || seen["b"] == 0 {
		t.Fatalf("seen=%v", seen)
	}
}
