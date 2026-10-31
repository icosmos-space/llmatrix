# LLMatrix Proxy

HTTP proxy service that fronts the Engine with auth, routing, rate limiting,
usage tracking, and billing.

## Run

```bash
# upstream keys (same env vars as SDK)
export OPENAI_API_KEY=...
# optional proxy auth for clients
export LLMATRIX_API_KEY=dev-key

go run ./cmd/llmatrix-proxy
```

OpenAI-compatible endpoints:

- `POST /v1/chat/completions` — model supports `provider/model` (e.g. `openai/gpt-4o-mini`)
- `POST /v1/embeddings`
- `GET /v1/models`
- `GET /healthz`

Proxy delegates to Engine providers; credentials come from provider env vars.

## Layout

| Path | Role |
|---|---|
| `cmd/llmatrix-proxy` | Process entrypoint |
| `internal/bootstrap` | Application wiring |
| `internal/server` | HTTP server |
| `internal/handler` | HTTP API handlers |
| `internal/middleware` | Auth, logging, rate limit, recovery |
| `internal/service` | Business logic |
| `internal/router` | Channel load-balance / fallback |
| `internal/auth` | API key authentication |
| `internal/billing` | Pricing and cost calculation |
| `internal/repository` | Persistence interfaces |
| `internal/config` | Configuration loading |
| `migrations` | DB migrations |
| `configs` | Sample configs |
| `deployments` | Docker / Kubernetes manifests |
| `web` | Optional admin UI (future) |
