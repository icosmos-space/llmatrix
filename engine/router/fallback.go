package router

import (
	"strings"

	"github.com/icosmos-space/llmatrix/engine/model"
)

// FallbackStrategy tries the primary selection, then ordered fallbacks.
// Callers typically walk Candidates() themselves when executing.
type FallbackStrategy struct {
	Primary   Selection
	Fallbacks []Selection
}

// Select returns the primary selection.
func (s FallbackStrategy) Select(_ SelectionRequest) (Selection, error) {
	if s.Primary.Provider == "" && s.Primary.Model == "" {
		return Selection{}, ErrNoStrategy
	}
	return s.Primary, nil
}

// Candidates returns primary followed by fallbacks.
func (s FallbackStrategy) Candidates() []Selection {
	out := make([]Selection, 0, 1+len(s.Fallbacks))
	if s.Primary.Provider != "" || s.Primary.Model != "" {
		out = append(out, s.Primary)
	}
	out = append(out, s.Fallbacks...)
	return out
}

// ParseModelRef builds a Selection from a LiteLLM-style "provider/model" string.
func ParseModelRef(raw string) Selection {
	ref := model.Resolve(raw, "")
	return Selection{Provider: ref.Provider, Model: ref.Model}
}

// ParseModelRefs parses a list of model refs.
func ParseModelRefs(raws ...string) []Selection {
	out := make([]Selection, 0, len(raws))
	for _, raw := range raws {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		out = append(out, ParseModelRef(raw))
	}
	return out
}
