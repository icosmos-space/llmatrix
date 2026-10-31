package handler

import (
	"encoding/json"
	"net/http"

	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

func (h *Handler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req types.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	resp, err := h.chat.Complete(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.usage.Record(r.Context(), resp.Provider, resp.Model, resp.Usage); err != nil {
		// usage persistence failures should not fail the request
		_ = err
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadGateway
	if e, ok := engineerrors.As(err); ok {
		switch e.Kind {
		case engineerrors.KindInvalidRequest:
			status = http.StatusBadRequest
		case engineerrors.KindUnauthorized:
			status = http.StatusUnauthorized
		case engineerrors.KindNotFound:
			status = http.StatusNotFound
		case engineerrors.KindRateLimited:
			status = http.StatusTooManyRequests
		case engineerrors.KindTimeout:
			status = http.StatusGatewayTimeout
		}
	}
	http.Error(w, err.Error(), status)
}
