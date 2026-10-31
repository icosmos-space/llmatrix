package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/icosmos-space/llmatrix/proxy/internal/config"
	"github.com/icosmos-space/llmatrix/proxy/internal/handler"
	"github.com/icosmos-space/llmatrix/proxy/internal/middleware"
	"github.com/icosmos-space/llmatrix/proxy/internal/server"
	"github.com/icosmos-space/llmatrix/proxy/internal/service"
)

// App is the assembled proxy process.
type App struct {
	cfg    *config.Config
	server *server.Server
}

// New wires dependencies and returns a runnable App.
func New(cfg *config.Config) (*App, error) {
	chatSvc := service.NewChatService()
	keySvc := service.NewKeyService()
	channelSvc := service.NewChannelService()
	usageSvc := service.NewUsageService()

	mux := http.NewServeMux()
	h := handler.New(chatSvc, keySvc, channelSvc, usageSvc)
	h.Register(mux)

	var httpHandler http.Handler = mux
	httpHandler = middleware.Recovery(httpHandler)
	httpHandler = middleware.Logging(httpHandler)
	httpHandler = middleware.RateLimit(httpHandler)
	httpHandler = middleware.Auth(cfg.APIKey, httpHandler)

	srv := server.New(cfg.ListenAddr, httpHandler)
	return &App{cfg: cfg, server: srv}, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	fmt.Printf("llmatrix-proxy listening on %s\n", a.cfg.ListenAddr)
	return a.server.ListenAndServe()
}
