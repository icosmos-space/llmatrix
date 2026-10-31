package openaicompatible

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/errors"
)

func mapHTTPError(provider string, status int, body []byte) error {
	msg := string(body)
	var payload struct {
		Error *wireError `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Error != nil && payload.Error.Message != "" {
		msg = payload.Error.Message
	}
	kind := errors.KindProvider
	switch status {
	case http.StatusBadRequest:
		kind = errors.KindInvalidRequest
	case http.StatusUnauthorized, http.StatusForbidden:
		kind = errors.KindUnauthorized
	case http.StatusNotFound:
		kind = errors.KindNotFound
	case http.StatusTooManyRequests:
		kind = errors.KindRateLimited
	case http.StatusRequestTimeout, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		kind = errors.KindTimeout
	}
	return errors.WrapStatus(kind, provider, fmt.Sprintf("HTTP %d: %s", status, truncate(msg, 500)), status, nil)
}
