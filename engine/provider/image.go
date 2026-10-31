package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// ImageProvider generates images from text prompts.
type ImageProvider interface {
	Provider
	ImageGeneration(ctx context.Context, req *types.ImageRequest) (*types.ImageResponse, error)
}

// ImageEditProvider edits images (OpenAI /images/edits, LiteLLM image_edit()).
type ImageEditProvider interface {
	Provider
	ImageEdit(ctx context.Context, req *types.ImageEditRequest) (*types.ImageResponse, error)
}
