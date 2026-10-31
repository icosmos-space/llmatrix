package engine

import "github.com/icosmos-space/llmatrix/engine/registry"

// Engine orchestrates provider selection and execution.
type Engine struct {
	registry *registry.Registry
	executor *Executor
}

// New creates an Engine bound to a provider registry.
func New(reg *registry.Registry) *Engine {
	return &Engine{
		registry: reg,
		executor: NewExecutor(reg),
	}
}

// Executor exposes the underlying call executor.
func (e *Engine) Executor() *Executor {
	return e.executor
}
