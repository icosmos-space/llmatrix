package usage

import (
	"strings"
	"unicode/utf8"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// TokenCounter estimates token counts without a model-specific tokenizer.
// This mirrors LiteLLM's lightweight fallback (≈ chars/4), not tiktoken accuracy.
type TokenCounter struct{}

// CountText estimates tokens for a single string.
func (TokenCounter) CountText(text string) int {
	if text == "" {
		return 0
	}
	// Prefer rune length; CJK tends to be denser than Latin, so use a slightly
	// higher estimate when non-ASCII dominates.
	n := utf8.RuneCountInString(text)
	ascii := 0
	for _, r := range text {
		if r < utf8.RuneSelf {
			ascii++
		}
	}
	if n > 0 && float64(ascii)/float64(n) < 0.5 {
		est := (n*10 + 6) / 7 // ~1.4 chars/token
		if est < 1 {
			return 1
		}
		return est
	}
	est := (n + 3) / 4
	if est < 1 {
		return 1
	}
	return est
}

// CountMessages estimates prompt tokens for chat messages (LiteLLM token_counter shape).
func (c TokenCounter) CountMessages(msgs []types.Message) int {
	total := 0
	for _, m := range msgs {
		total += 4 // role / framing overhead
		total += c.CountText(string(m.Role))
		total += c.CountText(m.Content)
		total += c.CountText(m.Name)
		total += c.CountText(m.ToolCallID)
		for _, part := range m.Parts {
			total += c.CountText(part.Text)
		}
		for _, tc := range m.ToolCalls {
			total += c.CountText(tc.Name) + c.CountText(tc.Arguments)
		}
	}
	if total > 0 {
		total += 2 // reply priming
	}
	return total
}

// CountRequest estimates tokens for a chat request (messages + tools schema).
func (c TokenCounter) CountRequest(req *types.ChatRequest) int {
	if req == nil {
		return 0
	}
	n := c.CountMessages(req.Messages)
	for _, t := range req.Tools {
		n += c.CountText(t.Function.Name)
		n += c.CountText(t.Function.Description)
		n += c.CountText(strings.Join(keysOf(t.Function.Parameters), ","))
	}
	return n
}

func keysOf(m map[string]any) []string {
	if m == nil {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
