package router

// Request describes a routing decision input.
type Request struct {
	Model string
}

// Result is the chosen channel.
type Result struct {
	ChannelID string
	Provider  string
	Model     string
}

// Router selects an upstream channel for a model.
type Router struct {
	balancer *LoadBalancer
	fallback *Fallback
}

// New creates a Router.
func New(balancer *LoadBalancer, fallback *Fallback) *Router {
	return &Router{balancer: balancer, fallback: fallback}
}

// Route picks a channel for the request.
func (r *Router) Route(req Request) (Result, error) {
	if r.balancer != nil {
		if res, err := r.balancer.Pick(req); err == nil {
			return res, nil
		}
	}
	if r.fallback != nil {
		return r.fallback.Next(req)
	}
	return Result{}, ErrNoChannel
}
