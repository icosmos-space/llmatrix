package openaicompatible

import (
	"context"
	"net/url"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

const assistantsBeta = "assistants=v2"

func (p *Provider) doAssistantsJSON(ctx context.Context, method, path string, body, dest any) error {
	req, err := p.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	req.Header.Set("OpenAI-Beta", assistantsBeta)
	resp, err := p.httpClient().Do(req)
	if err != nil {
		return errors.Wrap(errors.KindProvider, p.Info().Name, "http request failed", err)
	}
	defer resp.Body.Close()
	return p.decodeJSONResponse(resp, dest)
}

// CreateAssistant implements provider.AssistantsProvider.
func (p *Provider) CreateAssistant(ctx context.Context, req *types.AssistantCreateRequest) (*types.Assistant, error) {
	if req == nil || req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	var out types.Assistant
	if err := p.doAssistantsJSON(ctx, "POST", "/assistants", req, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// GetAssistant implements provider.AssistantsProvider.
func (p *Provider) GetAssistant(ctx context.Context, assistantID string) (*types.Assistant, error) {
	if assistantID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "assistant id is required", nil)
	}
	var out types.Assistant
	if err := p.doAssistantsJSON(ctx, "GET", "/assistants/"+url.PathEscape(assistantID), nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// DeleteAssistant implements provider.AssistantsProvider.
func (p *Provider) DeleteAssistant(ctx context.Context, assistantID string) error {
	if assistantID == "" {
		return errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "assistant id is required", nil)
	}
	return p.doAssistantsJSON(ctx, "DELETE", "/assistants/"+url.PathEscape(assistantID), nil, nil)
}

// CreateThread implements provider.AssistantsProvider.
func (p *Provider) CreateThread(ctx context.Context, req *types.ThreadCreateRequest) (*types.Thread, error) {
	if req == nil {
		req = &types.ThreadCreateRequest{}
	}
	var out types.Thread
	if err := p.doAssistantsJSON(ctx, "POST", "/threads", req, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// GetThread implements provider.AssistantsProvider.
func (p *Provider) GetThread(ctx context.Context, threadID string) (*types.Thread, error) {
	if threadID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "thread id is required", nil)
	}
	var out types.Thread
	if err := p.doAssistantsJSON(ctx, "GET", "/threads/"+url.PathEscape(threadID), nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// AddMessage implements provider.AssistantsProvider.
func (p *Provider) AddMessage(ctx context.Context, threadID string, msg types.ThreadMessage) (*types.ThreadMessageObject, error) {
	if threadID == "" || msg.Role == "" || msg.Content == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "thread id, role and content are required", nil)
	}
	var out types.ThreadMessageObject
	if err := p.doAssistantsJSON(ctx, "POST", "/threads/"+url.PathEscape(threadID)+"/messages", msg, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// ListMessages implements provider.AssistantsProvider.
func (p *Provider) ListMessages(ctx context.Context, threadID string) (*types.ThreadMessageList, error) {
	if threadID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "thread id is required", nil)
	}
	var out types.ThreadMessageList
	if err := p.doAssistantsJSON(ctx, "GET", "/threads/"+url.PathEscape(threadID)+"/messages", nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// CreateRun implements provider.AssistantsProvider.
func (p *Provider) CreateRun(ctx context.Context, threadID string, req *types.RunCreateRequest) (*types.Run, error) {
	if threadID == "" || req == nil || req.AssistantID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "thread id and assistant_id are required", nil)
	}
	var out types.Run
	if err := p.doAssistantsJSON(ctx, "POST", "/threads/"+url.PathEscape(threadID)+"/runs", req, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// GetRun implements provider.AssistantsProvider.
func (p *Provider) GetRun(ctx context.Context, threadID, runID string) (*types.Run, error) {
	if threadID == "" || runID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "thread id and run id are required", nil)
	}
	var out types.Run
	path := "/threads/" + url.PathEscape(threadID) + "/runs/" + url.PathEscape(runID)
	if err := p.doAssistantsJSON(ctx, "GET", path, nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}
