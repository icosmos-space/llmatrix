package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Rerank implements provider.RerankProvider (Cohere / LiteLLM shaped /v1/rerank).
func (p *Provider) Rerank(ctx context.Context, req *types.RerankRequest) (*types.RerankResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" || req.Query == "" || len(req.Documents) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model, query, and documents are required", nil)
	}
	body := map[string]any{
		"model":     req.Model,
		"query":     req.Query,
		"documents": req.Documents,
	}
	if req.TopN != nil {
		body["top_n"] = *req.TopN
	}
	if req.ReturnDocuments != nil {
		body["return_documents"] = *req.ReturnDocuments
	}

	var raw struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Results []struct {
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relevance_score"`
			Document       any     `json:"document"`
		} `json:"results"`
		Meta struct {
			BilledUnits struct {
				SearchUnits int `json:"search_units"`
			} `json:"billed_units"`
		} `json:"meta"`
		Error *wireError `json:"error"`
	}
	// Try /rerank then /v1/rerank-compatible path under base (base already includes /v1).
	if err := p.doJSON(ctx, "POST", "/rerank", body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, raw.Error.Message, nil)
	}
	out := &types.RerankResponse{
		ID:       raw.ID,
		Model:    firstNonEmpty(raw.Model, req.Model),
		Provider: p.Info().Name,
		Results:  make([]types.RerankResult, 0, len(raw.Results)),
	}
	for _, r := range raw.Results {
		item := types.RerankResult{
			Index:          r.Index,
			RelevanceScore: r.RelevanceScore,
		}
		switch d := r.Document.(type) {
		case string:
			item.Document = d
		case map[string]any:
			if t, ok := d["text"].(string); ok {
				item.Document = t
			}
		}
		out.Results = append(out.Results, item)
	}
	return out, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
