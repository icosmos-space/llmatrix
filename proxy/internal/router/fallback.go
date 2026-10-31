package router

// Fallback walks a prioritized channel list.
type Fallback struct {
	Channels []Result
}

// Next returns the first configured fallback channel.
func (f *Fallback) Next(_ Request) (Result, error) {
	if len(f.Channels) == 0 {
		return Result{}, ErrNoChannel
	}
	return f.Channels[0], nil
}
