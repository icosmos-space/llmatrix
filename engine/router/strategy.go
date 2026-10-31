package router

import "errors"

// Strategy decides how to pick a provider/model.
type Strategy interface {
	Select(req SelectionRequest) (Selection, error)
}

// ErrNoStrategy is returned when no strategy is configured.
var ErrNoStrategy = errors.New("router: no strategy configured")

// FixedStrategy always returns the same provider and model.
type FixedStrategy struct {
	Provider string
	Model    string
}

// Select implements Strategy.
func (s FixedStrategy) Select(_ SelectionRequest) (Selection, error) {
	return Selection{Provider: s.Provider, Model: s.Model}, nil
}
