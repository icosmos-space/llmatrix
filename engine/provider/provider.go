package provider

// Info describes a registered provider.
type Info struct {
	Name         string
	DisplayName  string
	Capabilities Capabilities
}

// Capabilities advertises supported features.
type Capabilities struct {
	Chat       bool
	Streaming  bool
	Tools      bool
	Embeddings bool
}

// Provider is the root provider identity.
type Provider interface {
	Info() Info
}
