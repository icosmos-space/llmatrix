package bedrock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/icosmos-space/llmatrix/engine/aws"
	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatCompletion implements provider.ChatProvider via Bedrock Converse.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	if req == nil || req.Model == "" || len(req.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "bedrock", "model and messages are required", nil)
	}
	body, err := toConverseRequest(req)
	if err != nil {
		return nil, err
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, "bedrock", "marshal request", err)
	}
	endpoint := p.runtimeEndpoint(req.Model, "converse")
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(rawBody))
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "bedrock", "create request", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	creds := p.credentials()
	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
		return nil, errors.Wrap(errors.KindUnauthorized, "bedrock", "AWS credentials are required", nil)
	}
	if err := aws.SignRequest(httpReq, rawBody, creds, p.region(), "bedrock", time.Time{}); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "bedrock", "sign request", err)
	}
	resp, err := p.httpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, "bedrock", "http request failed", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, "bedrock", "read response", err)
	}
	if resp.StatusCode >= 400 {
		return nil, mapHTTPError(resp.StatusCode, data)
	}
	var raw converseResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errors.Wrap(errors.KindProvider, "bedrock", "decode response", err)
	}
	return fromConverseResponse(req.Model, &raw), nil
}
