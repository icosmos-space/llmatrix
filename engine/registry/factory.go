package registry

import "github.com/icosmos-space/llmatrix/engine/provider"

// Factory builds a Provider from configuration.
type Factory func(cfg map[string]string) (provider.Provider, error)

// FactoryRegistry maps provider names to constructors.
type FactoryRegistry struct {
	factories map[string]Factory
}

// NewFactoryRegistry creates an empty factory map.
func NewFactoryRegistry() *FactoryRegistry {
	return &FactoryRegistry{factories: make(map[string]Factory)}
}

// Register associates a name with a Factory.
func (f *FactoryRegistry) Register(name string, factory Factory) {
	f.factories[name] = factory
}

// Create instantiates a provider by name.
func (f *FactoryRegistry) Create(name string, cfg map[string]string) (provider.Provider, error) {
	factory, ok := f.factories[name]
	if !ok {
		return nil, ErrUnknownProvider(name)
	}
	return factory(cfg)
}

// ErrUnknownProvider reports a missing factory.
func ErrUnknownProvider(name string) error {
	return &UnknownProviderError{Name: name}
}

// UnknownProviderError is returned when no factory is registered.
type UnknownProviderError struct {
	Name string
}

func (e *UnknownProviderError) Error() string {
	return "registry: unknown provider " + e.Name
}
