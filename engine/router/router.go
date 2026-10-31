package router

// Router selects a provider/model for an incoming request.
type Router interface {
	Select(req SelectionRequest) (Selection, error)
}

// DefaultRouter applies a Strategy.
type DefaultRouter struct {
	Strategy Strategy
}

// Select delegates to the configured Strategy.
func (r *DefaultRouter) Select(req SelectionRequest) (Selection, error) {
	if r.Strategy == nil {
		return Selection{}, ErrNoStrategy
	}
	return r.Strategy.Select(req)
}
