package types

// FileCreateRequest uploads a file (OpenAI /files, LiteLLM create_file()).
type FileCreateRequest struct {
	File     []byte `json:"-"`
	FileName string `json:"-"`
	Purpose  string `json:"purpose"` // assistants | batch | fine-tune | ...
}

// FileObject is an OpenAI file metadata object.
type FileObject struct {
	ID        string `json:"id"`
	Object    string `json:"object,omitempty"`
	Bytes     int64  `json:"bytes,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	Filename  string `json:"filename,omitempty"`
	Purpose   string `json:"purpose,omitempty"`
	Status    string `json:"status,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

// FileListResponse lists files.
type FileListResponse struct {
	Object   string       `json:"object,omitempty"`
	Data     []FileObject `json:"data"`
	Provider string       `json:"provider,omitempty"`
}
