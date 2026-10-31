package router

import "errors"

// ErrNoChannel indicates no upstream channel is available.
var ErrNoChannel = errors.New("router: no channel available")

// LoadBalancer picks among healthy channels.
type LoadBalancer struct {
	Channels []Result
	next     int
}

// Pick returns the next channel in round-robin order.
func (lb *LoadBalancer) Pick(_ Request) (Result, error) {
	if len(lb.Channels) == 0 {
		return Result{}, ErrNoChannel
	}
	res := lb.Channels[lb.next%len(lb.Channels)]
	lb.next++
	return res, nil
}
