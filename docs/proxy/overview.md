# Proxy

Module path: `github.com/icosmos-space/llmatrix/proxy`

Binary: `proxy/cmd/llmatrix-proxy`.

Exposes OpenAI-style HTTP endpoints (`/v1/chat/completions`, `/v1/embeddings`,
`/v1/models`) and layers auth, rate limiting, channel routing, usage, and billing
on top of Engine.
