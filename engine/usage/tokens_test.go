package usage

import (
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestTokenCounter(t *testing.T) {
	var c TokenCounter
	if n := c.CountText("hello world"); n < 1 {
		t.Fatalf("latin=%d", n)
	}
	if n := c.CountText("你好世界"); n < 1 {
		t.Fatalf("cjk=%d", n)
	}
	n := c.CountMessages([]types.Message{
		{Role: types.RoleUser, Content: "hi"},
		{Role: types.RoleAssistant, Content: "hello"},
	})
	if n < 5 {
		t.Fatalf("messages=%d", n)
	}
}
