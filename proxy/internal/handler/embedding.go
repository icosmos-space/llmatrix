package handler

import (
	"encoding/json"
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func (h *Handler) Embeddings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req types.EmbeddingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	resp, err := h.chat.Embed(r.Context(), &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.usage.Record(r.Context(), resp.Provider, resp.Model, resp.Usage); err != nil {
		_ = err
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
