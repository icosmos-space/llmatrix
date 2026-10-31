# Proxy

Module path: `github.com/icosmos-space/llmatrix/proxy`

Binary: `proxy/cmd/llmatrix-proxy`.

## Endpoints

| Method | Path | Behavior |
|---|---|---|
| POST | `/v1/chat/completions` | Delegates to Engine providers (`provider/model`) |
| POST | `/v1/embeddings` | Delegates to Engine embedding providers |
| GET | `/v1/models` | Lists known model identifiers |
| GET | `/healthz` | Liveness |

Upstream credentials use the same env vars as the SDK (`OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, …).
Optional client auth: set `LLMATRIX_API_KEY` and send `Authorization: Bearer …`.
