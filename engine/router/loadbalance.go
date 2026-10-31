package router

import (
	"sync"
	"sync/atomic"
)

// Deployment is one routable upstream (LiteLLM model_list entry).
type Deployment struct {
	Name     string // logical model name clients request
	Provider string
	Model    string
	APIKey   string
	BaseURL  string
	Weight   int
}

// SimpleLoadBalance picks deployments for a logical model name (round-robin).
type SimpleLoadBalance struct {
	mu       sync.RWMutex
	byName   map[string][]Deployment
	counters map[string]*atomic.Uint64
}

// NewSimpleLoadBalance creates an empty load balancer.
func NewSimpleLoadBalance() *SimpleLoadBalance {
	return &SimpleLoadBalance{
		byName:   make(map[string][]Deployment),
		counters: make(map[string]*atomic.Uint64),
	}
}

// Register adds a deployment under its logical Name.
func (lb *SimpleLoadBalance) Register(d Deployment) {
	if d.Weight <= 0 {
		d.Weight = 1
	}
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.byName[d.Name] = append(lb.byName[d.Name], d)
	if _, ok := lb.counters[d.Name]; !ok {
		lb.counters[d.Name] = &atomic.Uint64{}
	}
}

// Select implements Strategy using ModelHint as the logical model name.
func (lb *SimpleLoadBalance) Select(req SelectionRequest) (Selection, error) {
	d, err := lb.Pick(req.ModelHint)
	if err != nil {
		return Selection{}, err
	}
	return Selection{Provider: d.Provider, Model: d.Model}, nil
}

// Pick returns the next deployment for a logical model name.
func (lb *SimpleLoadBalance) Pick(name string) (Deployment, error) {
	lb.mu.RLock()
	deps := lb.byName[name]
	counter := lb.counters[name]
	lb.mu.RUnlock()
	if len(deps) == 0 {
		return Deployment{}, ErrNoChannel
	}
	// Weighted round-robin via expanded slots.
	slots := make([]int, 0, len(deps))
	for i, d := range deps {
		for w := 0; w < d.Weight; w++ {
			slots = append(slots, i)
		}
	}
	idx := int(counter.Add(1)-1) % len(slots)
	return deps[slots[idx]], nil
}

// ErrNoChannel indicates no deployment is registered for the model.
var ErrNoChannel = errNoChannel

type errNoChannelT struct{}

func (errNoChannelT) Error() string { return "router: no deployment for model" }

var errNoChannel error = errNoChannelT{}
