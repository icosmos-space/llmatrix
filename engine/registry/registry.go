package registry

import (
	"fmt"
	"sync"

	"github.com/icosmos-space/llmatrix/engine/provider"
)

// Registry stores named providers.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]provider.Provider
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{providers: make(map[string]provider.Provider)}
}

// Register adds a provider under its Info().Name.
func (r *Registry) Register(p provider.Provider) error {
	info := p.Info()
	if info.Name == "" {
		return fmt.Errorf("registry: provider name is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.providers[info.Name]; exists {
		return fmt.Errorf("registry: provider %q already registered", info.Name)
	}
	r.providers[info.Name] = p
	return nil
}

// Get returns a provider by name.
func (r *Registry) Get(name string) (provider.Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	return p, ok
}

// List returns all registered provider names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
