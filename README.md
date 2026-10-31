# LLMatrix

<p align="center">
  <strong>Unified LLM SDK for Go</strong>
  <br />
  One API. Multiple Providers. Go Native.
</p>

<p align="center">
  <a href="https://github.com/icosmos-space/llmatrix/stargazers"><img src="https://img.shields.io/github/stars/icosmos-space/llmatrix?style=flat" alt="GitHub Stars" /></a>
  <a href="https://github.com/icosmos-space/llmatrix/blob/main/LICENSE"><img src="https://img.shields.io/github/license/icosmos-space/llmatrix" alt="License" /></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go" alt="Go Version" /></a>
  <a href="https://github.com/icosmos-space/llmatrix/actions"><img src="https://img.shields.io/github/actions/workflow/status/icosmos-space/llmatrix/ci.yml?branch=main" alt="Build Status" /></a>
</p>

**LLMatrix** is a unified LLM SDK written in Go. It provides a consistent interface for interacting with multiple Large Language Model providers, allowing developers to switch between models without rewriting application code.

Inspired by the unified API philosophy of [LiteLLM](https://github.com/BerriAI/litellm), LLMatrix aims to bring a lightweight, idiomatic, and extensible LLM integration experience to the Go ecosystem.

## Why LLMatrix?

Integrating multiple LLM providers often means dealing with different APIs, request formats, response structures, authentication methods, and streaming protocols.

LLMatrix abstracts these differences behind a unified Go interface.

- **Unified API** — Use one interface across multiple LLM providers.
- **Provider Agnostic** — Switch providers with minimal code changes.
- **Go Native** — Designed with Go conventions, contexts, interfaces, and concurrency in mind.
- **Streaming First** — Support incremental response generation.
- **Extensible Architecture** — Add new providers without changing application-level code.
- **Lightweight** — Keep dependencies and abstractions focused.

## Supported Providers

| Provider | Chat Completion | Streaming | Tool Calling | Embeddings |
|---|---|---|---|---|
| OpenAI | Yes | Yes | Yes | Yes |
| DeepSeek | Yes | Yes | Yes | Yes |
| Groq | Yes | Yes | Yes | Yes |
| Together | Yes | Yes | Yes | Yes |
| Mistral | Yes | Yes | Yes | Yes |
| OpenRouter | Yes | Yes | Yes | Yes |
| xAI | Yes | Yes | Yes | Yes |
| Fireworks | Yes | Yes | Yes | Yes |
| Ollama | Yes | Yes | Yes | Yes |
| OpenAI-compatible APIs | Yes | Yes | Yes | Yes |
| Anthropic | Yes | Yes | Yes | — |
| Google Gemini | Yes | Yes | Yes | Yes |
| Azure OpenAI | Yes | Yes | Yes | Yes |

> Provider support is under active development. Check the project roadmap for implementation status.

## Repository layout

```text
llmatrix/
├── sdk/        # Public Go client
├── engine/     # Provider execution core
├── proxy/      # HTTP proxy service
├── examples/   # Cross-module examples
├── tests/      # Integration / e2e
└── docs/       # Architecture and module docs
```

## Installation

```bash
go get github.com/icosmos-space/llmatrix/sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/icosmos-space/llmatrix/sdk"
    "github.com/icosmos-space/llmatrix/sdk/chat"
)

func main() {
    client := sdk.NewClient(
        sdk.WithProvider("openai"),
        sdk.WithAPIKey("YOUR_API_KEY"),
        sdk.WithModel("gpt-4o"),
    )

    resp, err := client.ChatCompletion(
        context.Background(),
        &chat.Request{
            Messages: []chat.Message{
                {
                    Role:    "user",
                    Content: "Explain how Go interfaces work.",
                },
            },
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
```

> Supported: OpenAI / DeepSeek / Ollama / Anthropic / Gemini (and any OpenAI-compatible base URL). Non-streaming calls include default retries.

## Architecture

LLMatrix is a Go workspace with three modules: **SDK**, **Engine**, and **Proxy**.

```text
Application / SDK Client
        |
        v
   +---------+
   |   SDK   |
   +---------+
        |
        +------------------+
        |                  |
        v                  v
   +---------+       +-----------+
   | Engine  | <---- |   Proxy   |
   +---------+       +-----------+
        |
   providers / protocol / router
```

| Module | Path | Responsibility |
|---|---|---|
| SDK | `sdk/` | Public client API (chat, embedding, streaming, tools, proxy client) |
| Engine | `engine/` | Unified types, provider interfaces, registry, execution, routing |
| Proxy | `proxy/` | HTTP gateway with auth, rate limit, channel routing, billing |

See [docs/architecture/overview.md](docs/architecture/overview.md) for details.

## Roadmap

### Phase 1 — Core SDK

- [x] Unified client interface
- [x] Provider adapter architecture
- [x] OpenAI-compatible API support
- [x] Chat completion
- [x] Streaming responses
- [x] Context cancellation and timeout
- [x] Unified error handling (HTTP status → engine error kinds)

### Phase 2 — Provider Ecosystem

- [x] OpenAI
- [x] Anthropic
- [x] Google Gemini
- [x] DeepSeek
- [x] Ollama
- [ ] Additional providers (Azure, Bedrock, …)

### Phase 3 — Advanced Features

- [x] Tool Calling (OpenAI-compatible + Anthropic)
- [x] Embeddings (OpenAI-compatible + Gemini)
- [x] Model aliases
- [x] Automatic retries
- [x] Fallback strategies
- [x] Load balancing (`sdk.Router`)
- [x] Model routing (alias + provider/model)
- [x] Proxy → Engine chat/embeddings wiring
- [x] Image / speech / transcription / rerank (OpenAI-compatible)
- [x] OpenAI Responses API (`responses()` + chat bridge)
- [x] Azure OpenAI provider
- [x] Gemini tool calling

### Phase 4 — Developer Experience

- [ ] Structured logging
- [ ] Usage tracking
- [ ] Token usage statistics
- [ ] Middleware support
- [ ] Observability hooks
- [ ] Comprehensive examples
- [ ] Benchmarks

## Design Principles

LLMatrix follows a few core principles:

1. **Unified, not restrictive** — Normalize common functionality while preserving provider-specific capabilities when needed.
2. **Explicit over implicit** — Make configuration, errors, and provider behavior predictable.
3. **Composable architecture** — Keep providers, routing, retry policies, and middleware independently extensible.
4. **Go idiomatic** — Prefer standard Go patterns over unnecessary abstractions.
5. **No unnecessary dependencies** — Keep the core SDK small and focused.

## Contributing

Contributions, bug reports, feature requests, and provider integrations are welcome.

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request.

Please open an issue before implementing major architectural changes.

## License

LLMatrix is licensed under the [MIT License](LICENSE).

---

<p align="center">
  Built with Go ❤️ for the LLM ecosystem.
</p>
