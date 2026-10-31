package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/assistants"
)

func (c *Client) assistantsProvider(modelHint string) (provider.AssistantsProvider, error) {
	raw := modelHint
	if raw == "" {
		raw = c.cfg.Model
	}
	if raw == "" {
		prov := c.cfg.Provider
		if prov == "" {
			prov = "openai"
		}
		raw = prov + "/gpt-4o"
	}
	refs := c.candidates(raw)
	if len(refs) == 0 {
		return nil, fmt.Errorf("%w: model is required", ErrInvalidConfig)
	}
	p, err := c.newProvider(refs[0], true)
	if err != nil {
		return nil, err
	}
	ap, ok := p.(provider.AssistantsProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support assistants", ErrUnsupported, refs[0].Provider)
	}
	return ap, nil
}

// CreateAssistants is LiteLLM create_assistants().
func (c *Client) CreateAssistants(ctx context.Context, req *assistants.CreateRequest) (*assistants.Assistant, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	ap, err := c.assistantsProvider(req.Model)
	if err != nil {
		return nil, err
	}
	modelName := req.Model
	if refs := c.candidates(req.Model); len(refs) > 0 {
		modelName = refs[0].Model
	}
	out, err := ap.CreateAssistant(ctx, &types.AssistantCreateRequest{
		Model: modelName, Name: req.Name, Description: req.Description,
		Instructions: req.Instructions, Tools: req.Tools, Metadata: req.Metadata, Temperature: req.Temperature,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineAssistant(out), nil
}

// GetAssistant fetches an assistant.
func (c *Client) GetAssistant(ctx context.Context, assistantID string, modelHint ...string) (*assistants.Assistant, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := ap.GetAssistant(ctx, assistantID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineAssistant(out), nil
}

// DeleteAssistant deletes an assistant.
func (c *Client) DeleteAssistant(ctx context.Context, assistantID string, modelHint ...string) error {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return err
	}
	if err := ap.DeleteAssistant(ctx, assistantID); err != nil {
		return MapError(err)
	}
	return nil
}

// CreateThread is LiteLLM create_thread().
func (c *Client) CreateThread(ctx context.Context, req *assistants.ThreadCreateRequest) (*assistants.Thread, error) {
	if req == nil {
		req = &assistants.ThreadCreateRequest{}
	}
	ap, err := c.assistantsProvider(req.Model)
	if err != nil {
		return nil, err
	}
	engineReq := &types.ThreadCreateRequest{Metadata: req.Metadata}
	for _, m := range req.Messages {
		engineReq.Messages = append(engineReq.Messages, types.ThreadMessage{Role: m.Role, Content: m.Content})
	}
	out, err := ap.CreateThread(ctx, engineReq)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineThread(out), nil
}

// GetThread is LiteLLM get_thread().
func (c *Client) GetThread(ctx context.Context, threadID string, modelHint ...string) (*assistants.Thread, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := ap.GetThread(ctx, threadID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineThread(out), nil
}

// AddMessage is LiteLLM add_message().
func (c *Client) AddMessage(ctx context.Context, threadID string, msg assistants.Message, modelHint ...string) (*assistants.MessageObject, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := ap.AddMessage(ctx, threadID, types.ThreadMessage{Role: msg.Role, Content: msg.Content})
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineThreadMessage(out), nil
}

// ListMessages lists messages on a thread.
func (c *Client) ListMessages(ctx context.Context, threadID string, modelHint ...string) (*assistants.MessageList, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := ap.ListMessages(ctx, threadID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineThreadMessageList(out), nil
}

// CreateRun starts a thread run.
func (c *Client) CreateRun(ctx context.Context, threadID string, req *assistants.RunCreateRequest, modelHint ...string) (*assistants.Run, error) {
	if req == nil || req.AssistantID == "" {
		return nil, ErrInvalidRequest
	}
	hint := req.Model
	if len(modelHint) > 0 && modelHint[0] != "" {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := ap.CreateRun(ctx, threadID, &types.RunCreateRequest{
		AssistantID: req.AssistantID, Model: req.Model, Instructions: req.Instructions,
		Tools: req.Tools, Metadata: req.Metadata, Temperature: req.Temperature,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineRun(out), nil
}

// GetRun fetches a run.
func (c *Client) GetRun(ctx context.Context, threadID, runID string, modelHint ...string) (*assistants.Run, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	ap, err := c.assistantsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := ap.GetRun(ctx, threadID, runID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineRun(out), nil
}

func fromEngineAssistant(a *types.Assistant) *assistants.Assistant {
	if a == nil {
		return nil
	}
	return &assistants.Assistant{
		ID: a.ID, Object: a.Object, CreatedAt: a.CreatedAt, Name: a.Name, Description: a.Description,
		Model: a.Model, Instructions: a.Instructions, Tools: a.Tools, Metadata: a.Metadata, Provider: a.Provider,
	}
}

func fromEngineThread(t *types.Thread) *assistants.Thread {
	if t == nil {
		return nil
	}
	return &assistants.Thread{ID: t.ID, Object: t.Object, CreatedAt: t.CreatedAt, Metadata: t.Metadata, Provider: t.Provider}
}

func fromEngineThreadMessage(m *types.ThreadMessageObject) *assistants.MessageObject {
	if m == nil {
		return nil
	}
	return &assistants.MessageObject{
		ID: m.ID, Object: m.Object, CreatedAt: m.CreatedAt, ThreadID: m.ThreadID, Role: m.Role,
		Content: m.Content, AssistantID: m.AssistantID, RunID: m.RunID, Provider: m.Provider,
	}
}

func fromEngineThreadMessageList(l *types.ThreadMessageList) *assistants.MessageList {
	if l == nil {
		return nil
	}
	out := &assistants.MessageList{Object: l.Object, Provider: l.Provider}
	for i := range l.Data {
		out.Data = append(out.Data, *fromEngineThreadMessage(&l.Data[i]))
	}
	return out
}

func fromEngineRun(r *types.Run) *assistants.Run {
	if r == nil {
		return nil
	}
	return &assistants.Run{
		ID: r.ID, Object: r.Object, CreatedAt: r.CreatedAt, ThreadID: r.ThreadID, AssistantID: r.AssistantID,
		Status: r.Status, Model: r.Model, Instructions: r.Instructions, Provider: r.Provider,
	}
}
