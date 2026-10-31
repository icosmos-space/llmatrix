package files

// CreateRequest uploads a file (LiteLLM create_file()).
type CreateRequest struct {
	File     []byte `json:"-"`
	FileName string `json:"-"`
	Purpose  string `json:"purpose"`
	// Model routes credentials via provider/model when set (LiteLLM-style).
	Model string `json:"model,omitempty"`
}

// Object is an OpenAI file metadata object.
type Object struct {
	ID        string `json:"id"`
	Object    string `json:"object,omitempty"`
	Bytes     int64  `json:"bytes,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	Filename  string `json:"filename,omitempty"`
	Purpose   string `json:"purpose,omitempty"`
	Status    string `json:"status,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

// ListResponse lists files.
type ListResponse struct {
	Object   string   `json:"object,omitempty"`
	Data     []Object `json:"data"`
	Provider string   `json:"provider,omitempty"`
}
