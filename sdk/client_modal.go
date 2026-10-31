package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/audio"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
	"github.com/icosmos-space/llmatrix/sdk/image"
	"github.com/icosmos-space/llmatrix/sdk/moderation"
	"github.com/icosmos-space/llmatrix/sdk/rerank"
)

// Embeddings creates embeddings (LiteLLM embedding()).
func (c *Client) Embeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	modelKey := c.cacheModelKey(req.Model)
	if c.cacheEnabled(req.Caching) {
		var cached embedding.Response
		if c.cacheGet("embedding", modelKey, req, &cached) {
			return &cached, nil
		}
	}
	out, err := c.withCandidates(ctx, "embedding", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		embProv, ok := p.(provider.EmbeddingProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support embeddings", ErrUnsupported, ref.Provider)
		}
		resp, err := embProv.Embeddings(ctx, toEngineEmbeddingRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return resp.Usage, fromEngineEmbeddingResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	resp := out.(*embedding.Response)
	if c.cacheEnabled(req.Caching) {
		c.cacheSet("embedding", modelKey, req, resp)
	}
	return resp, nil
}

// ImageGeneration creates images from a prompt (LiteLLM image_generation()).
func (c *Client) ImageGeneration(ctx context.Context, req *image.Request) (*image.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	out, err := c.withCandidates(ctx, "image", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		imgProv, ok := p.(provider.ImageProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support image generation", ErrUnsupported, ref.Provider)
		}
		resp, err := imgProv.ImageGeneration(ctx, toEngineImageRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return types.Usage{}, fromEngineImageResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*image.Response), nil
}

// Speech synthesizes audio (LiteLLM speech()).
func (c *Client) Speech(ctx context.Context, req *audio.SpeechRequest) (*audio.SpeechResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	out, err := c.withCandidates(ctx, "speech", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		sp, ok := p.(provider.SpeechProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support speech", ErrUnsupported, ref.Provider)
		}
		resp, err := sp.Speech(ctx, toEngineSpeechRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return types.Usage{}, fromEngineSpeechResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*audio.SpeechResponse), nil
}

// Transcription transcribes audio (LiteLLM transcription()).
func (c *Client) Transcription(ctx context.Context, req *audio.TranscriptionRequest) (*audio.TranscriptionResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	out, err := c.withCandidates(ctx, "transcription", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		tp, ok := p.(provider.TranscriptionProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support transcription", ErrUnsupported, ref.Provider)
		}
		resp, err := tp.Transcription(ctx, toEngineTranscriptionRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return types.Usage{}, fromEngineTranscriptionResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*audio.TranscriptionResponse), nil
}

// Rerank ranks documents (LiteLLM rerank()).
func (c *Client) Rerank(ctx context.Context, req *rerank.Request) (*rerank.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	out, err := c.withCandidates(ctx, "rerank", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		rp, ok := p.(provider.RerankProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support rerank", ErrUnsupported, ref.Provider)
		}
		resp, err := rp.Rerank(ctx, toEngineRerankRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return resp.Usage, fromEngineRerankResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*rerank.Response), nil
}

// ImageVariation creates image variations.
func (c *Client) ImageVariation(ctx context.Context, req *image.VariationRequest) (*image.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	modelName := req.Model
	if modelName == "" {
		modelName = "openai/dall-e-2"
	}
	out, err := c.withCandidates(ctx, "image_variation", modelName, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		vp, ok := p.(provider.ImageVariationProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support image variation", ErrUnsupported, ref.Provider)
		}
		resp, err := vp.ImageVariation(ctx, &types.ImageVariationRequest{
			Model: ref.Model, Image: req.Image, ImageFileName: req.ImageFileName,
			N: req.N, Size: req.Size, ResponseFormat: req.ResponseFormat, User: req.User,
		})
		if err != nil {
			return types.Usage{}, nil, err
		}
		return types.Usage{}, fromEngineImageResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*image.Response), nil
}

// ImageEdit edits an image (LiteLLM image_edit()).
func (c *Client) ImageEdit(ctx context.Context, req *image.EditRequest) (*image.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	modelName := req.Model
	if modelName == "" {
		modelName = "openai/dall-e-2"
	}
	out, err := c.withCandidates(ctx, "image_edit", modelName, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		ep, ok := p.(provider.ImageEditProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support image edit", ErrUnsupported, ref.Provider)
		}
		resp, err := ep.ImageEdit(ctx, toEngineImageEditRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return types.Usage{}, fromEngineImageResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*image.Response), nil
}

// Moderation classifies content (LiteLLM moderation()).
func (c *Client) Moderation(ctx context.Context, req *moderation.Request) (*moderation.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	modelName := req.Model
	if modelName == "" {
		modelName = "openai/omni-moderation-latest"
	}
	out, err := c.withCandidates(ctx, "moderation", modelName, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		mp, ok := p.(provider.ModerationProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support moderation", ErrUnsupported, ref.Provider)
		}
		resp, err := mp.Moderation(ctx, toEngineModerationRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return types.Usage{}, fromEngineModerationResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*moderation.Response), nil
}
