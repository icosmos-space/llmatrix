package router

// SelectionRequest describes routing inputs.
type SelectionRequest struct {
	ModelHint    string
	ProviderHint string
	Tags         []string
}

// Selection is the chosen provider and model.
type Selection struct {
	Provider string
	Model    string
}
