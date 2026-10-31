package sdk

import (
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/audio"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
	"github.com/icosmos-space/llmatrix/sdk/image"
	"github.com/icosmos-space/llmatrix/sdk/rerank"
	"github.com/icosmos-space/llmatrix/sdk/responses"
	"github.com/icosmos-space/llmatrix/sdk/tools"
)

func toEngineChatRequest(req *chat.Request, model string) *types.ChatRequest {
	out := &types.ChatRequest{
		Model:             model,
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 req.N,
		Stop:              req.Stop,
		MaxTokens:         req.MaxTokens,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		LogitBias:         req.LogitBias,
		User:              req.User,
		Seed:              req.Seed,
		ResponseFormat:    req.ResponseFormat,
		Stream:            req.Stream,
		StreamOptions:     req.StreamOptions,
		ToolChoice:        req.ToolChoice,
		ParallelToolCalls: req.ParallelToolCalls,
		Messages:          make([]types.Message, 0, len(req.Messages)),
	}
	out.PromptCacheKey = req.PromptCacheKey
	for _, m := range req.Messages {
		msg := types.Message{
			Role:       types.Role(m.Role),
			Content:    m.Content,
			Name:       m.Name,
			ToolCallID: m.ToolCallID,
		}
		if m.CacheControl != nil {
			msg.CacheControl = &types.CacheControl{Type: m.CacheControl.Type}
		}
		for _, part := range m.Parts {
			cp := types.ContentPart{Type: part.Type, Text: part.Text}
			if part.CacheControl != nil {
				cp.CacheControl = &types.CacheControl{Type: part.CacheControl.Type}
			}
			msg.Parts = append(msg.Parts, cp)
		}
		for _, tc := range m.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, types.ToolCall{
				ID:        tc.ID,
				Type:      tc.Type,
				Name:      tc.Name,
				Arguments: tc.Arguments,
			})
		}
		out.Messages = append(out.Messages, msg)
	}
	for _, t := range req.Tools {
		out.Tools = append(out.Tools, types.Tool{
			Type: t.Type,
			Function: types.ToolFunction{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		})
	}
	return out
}

func fromEngineChatResponse(resp *types.ChatResponse) *chat.Response {
	if resp == nil {
		return nil
	}
	out := &chat.Response{
		ID:       resp.ID,
		Model:    resp.Model,
		Provider: resp.Provider,
		Usage: chat.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Choices: make([]chat.Choice, 0, len(resp.Choices)),
	}
	for _, c := range resp.Choices {
		msg := chat.Message{
			Role:       string(c.Message.Role),
			Content:    c.Message.Content,
			Name:       c.Message.Name,
			ToolCallID: c.Message.ToolCallID,
		}
		for _, tc := range c.Message.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, tools.Call{
				ID:        tc.ID,
				Type:      tc.Type,
				Name:      tc.Name,
				Arguments: tc.Arguments,
			})
		}
		out.Choices = append(out.Choices, chat.Choice{
			Index:        c.Index,
			Message:      msg,
			FinishReason: c.FinishReason,
		})
	}
	return out
}

func toEngineEmbeddingRequest(req *embedding.Request, model string) *types.EmbeddingRequest {
	return &types.EmbeddingRequest{
		Model:          model,
		Input:          req.Input,
		Dimensions:     req.Dimensions,
		EncodingFormat: req.EncodingFormat,
		User:           req.User,
	}
}

func fromEngineEmbeddingResponse(resp *types.EmbeddingResponse) *embedding.Response {
	if resp == nil {
		return nil
	}
	out := &embedding.Response{
		Model:    resp.Model,
		Provider: resp.Provider,
		Usage: embedding.Usage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
		Data: make([]embedding.Item, 0, len(resp.Data)),
	}
	for _, d := range resp.Data {
		out.Data = append(out.Data, embedding.Item{
			Index:     d.Index,
			Embedding: d.Embedding,
		})
	}
	return out
}

func toEngineImageRequest(req *image.Request, model string) *types.ImageRequest {
	return &types.ImageRequest{
		Model:          model,
		Prompt:         req.Prompt,
		N:              req.N,
		Size:           req.Size,
		Quality:        req.Quality,
		Style:          req.Style,
		ResponseFormat: req.ResponseFormat,
		User:           req.User,
	}
}

func fromEngineImageResponse(resp *types.ImageResponse) *image.Response {
	if resp == nil {
		return nil
	}
	out := &image.Response{
		Created:  resp.Created,
		Provider: resp.Provider,
		Model:    resp.Model,
		Data:     make([]image.Data, 0, len(resp.Data)),
	}
	for _, d := range resp.Data {
		out.Data = append(out.Data, image.Data{
			URL:     d.URL,
			B64JSON: d.B64JSON,
			Revised: d.Revised,
		})
	}
	return out
}

func toEngineSpeechRequest(req *audio.SpeechRequest, model string) *types.SpeechRequest {
	return &types.SpeechRequest{
		Model:          model,
		Input:          req.Input,
		Voice:          req.Voice,
		ResponseFormat: req.ResponseFormat,
		Speed:          req.Speed,
	}
}

func fromEngineSpeechResponse(resp *types.SpeechResponse) *audio.SpeechResponse {
	if resp == nil {
		return nil
	}
	return &audio.SpeechResponse{
		ContentType: resp.ContentType,
		Audio:       resp.Audio,
		Provider:    resp.Provider,
		Model:       resp.Model,
	}
}

func toEngineTranscriptionRequest(req *audio.TranscriptionRequest, model string) *types.TranscriptionRequest {
	return &types.TranscriptionRequest{
		Model:          model,
		FileName:       req.FileName,
		Audio:          req.Audio,
		Language:       req.Language,
		Prompt:         req.Prompt,
		ResponseFormat: req.ResponseFormat,
		Temperature:    req.Temperature,
	}
}

func fromEngineTranscriptionResponse(resp *types.TranscriptionResponse) *audio.TranscriptionResponse {
	if resp == nil {
		return nil
	}
	return &audio.TranscriptionResponse{
		Text:     resp.Text,
		Language: resp.Language,
		Duration: resp.Duration,
		Provider: resp.Provider,
		Model:    resp.Model,
	}
}

func toEngineRerankRequest(req *rerank.Request, model string) *types.RerankRequest {
	return &types.RerankRequest{
		Model:           model,
		Query:           req.Query,
		Documents:       req.Documents,
		TopN:            req.TopN,
		ReturnDocuments: req.ReturnDocuments,
	}
}

func toEngineResponsesRequest(req *responses.Request, model string) *types.ResponsesRequest {
	return &types.ResponsesRequest{
		Model:           model,
		Input:           req.Input,
		Instructions:    req.Instructions,
		MaxOutputTokens: req.MaxOutputTokens,
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		Tools:           req.Tools,
		ToolChoice:      req.ToolChoice,
		Stream:          req.Stream,
		User:            req.User,
		Metadata:        req.Metadata,
		Store:           req.Store,
	}
}

func fromEngineResponsesResponse(resp *types.ResponsesResponse) *responses.Response {
	if resp == nil {
		return nil
	}
	out := &responses.Response{
		ID:       resp.ID,
		Object:   resp.Object,
		Created:  resp.Created,
		Status:   resp.Status,
		Model:    resp.Model,
		Provider: resp.Provider,
		Usage: responses.Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
		Output: make([]responses.OutputItem, 0, len(resp.Output)),
	}
	for _, item := range resp.Output {
		oi := responses.OutputItem{
			Type:   item.Type,
			ID:     item.ID,
			Status: item.Status,
			Role:   item.Role,
			Result: item.Result,
		}
		for _, part := range item.Content {
			oi.Content = append(oi.Content, responses.ContentPart{Type: part.Type, Text: part.Text})
		}
		out.Output = append(out.Output, oi)
	}
	return out
}

func fromEngineRerankResponse(resp *types.RerankResponse) *rerank.Response {
	if resp == nil {
		return nil
	}
	out := &rerank.Response{
		ID:       resp.ID,
		Model:    resp.Model,
		Provider: resp.Provider,
		Results:  make([]rerank.Result, 0, len(resp.Results)),
		Usage: rerank.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}
	for _, r := range resp.Results {
		out.Results = append(out.Results, rerank.Result{
			Index:          r.Index,
			RelevanceScore: r.RelevanceScore,
			Document:       r.Document,
		})
	}
	return out
}
