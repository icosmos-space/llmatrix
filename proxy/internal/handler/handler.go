package handler

import (
	"net/http"

	"github.com/icosmos-space/llmatrix/proxy/internal/service"
)

// Handler groups HTTP endpoints.
type Handler struct {
	chat    *service.ChatService
	key     *service.KeyService
	channel *service.ChannelService
	usage   *service.UsageService
}

// New constructs a Handler.
func New(chat *service.ChatService, key *service.KeyService, channel *service.ChannelService, usage *service.UsageService) *Handler {
	return &Handler{chat: chat, key: key, channel: channel, usage: usage}
}

// Register mounts routes on mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/chat/completions", h.ChatCompletions)
	mux.HandleFunc("/v1/embeddings", h.Embeddings)
	mux.HandleFunc("/v1/models", h.ListModels)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
