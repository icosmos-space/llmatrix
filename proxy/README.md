# LLMatrix Proxy

HTTP proxy service that fronts the Engine with auth, routing, rate limiting,
usage tracking, and billing.

## Run (scaffold)

```bash
go run ./cmd/llmatrix-proxy
```

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
