# LLMatrix

<p align="center">
  <strong>Unified LLM SDK for Go</strong>
  <br />
  One API. Multiple Providers. Go Native.
</p>

<p align="center">
  <a href="https://github.com/your-username/llmatrix/stargazers"><img src="https://img.shields.io/github/stars/your-username/llmatrix?style=flat" alt="GitHub Stars" /></a>
  <a href="https://github.com/your-username/llmatrix/blob/main/LICENSE"><img src="https://img.shields.io/github/license/your-username/llmatrix" alt="License" /></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go" alt="Go Version" /></a>
  <a href="https://github.com/your-username/llmatrix/actions"><img src="https://img.shields.io/github/actions/workflow/status/your-username/llmatrix/ci.yml?branch=main" alt="Build Status" /></a>
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
| OpenAI | Planned | Planned | Planned | Planned |
| Anthropic | Planned | Planned | Planned | Planned |
| Google Gemini | Planned | Planned | Planned | Planned |
| DeepSeek | Planned | Planned | Planned | Planned |
| Ollama | Planned | Planned | Planned | Planned |
| OpenAI-compatible APIs | Planned | Planned | Planned | Planned |

> Provider support is under active development. Check the project roadmap for implementation status.

## Installation

```bash
go get github.com/your-username/llmatrix
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/your-username/llmatrix"
)

func main() {
    client := llmatrix.NewClient(
        llmatrix.WithProvider("openai"),
        llmatrix.WithAPIKey("YOUR_API_KEY"),
        llmatrix.WithModel("gpt-4o"),
    )

    resp, err := client.ChatCompletion(
        context.Background(),
        &llmatrix.ChatRequest{
            Messages: []llmatrix.Message{
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

> The API above illustrates the intended interface. It will be updated as the public API stabilizes.

## Architecture

LLMatrix separates the unified API from provider-specific implementations.

```text
                 Application
                      |
                      v
              +---------------+
              |    LLMatrix   |
              |  Unified API  |
              +---------------+
                      |
                      v
              +---------------+
              | Provider Layer|
              +---------------+
                /     |     \
               /      |      \
              v       v       v
           OpenAI  Anthropic  Gemini
              |       |       |
              v       v       v
           Model APIs and Services
```

### Core Modules

| Module | Responsibility |
|---|---|
| Client | Unified entry point |
| Provider | Provider-specific API adapters |
| Chat | Chat completion requests and responses |
| Streaming | Stream response handling |
| Embedding | Text embedding interfaces |
| Tool Calling | Function and tool invocation |
| Router | Model selection and routing |
| Retry | Retry and error recovery |
| Types | Shared request and response structures |

## Roadmap

### Phase 1 — Core SDK

- [ ] Unified client interface
- [ ] Provider adapter architecture
- [ ] OpenAI-compatible API support
- [ ] Chat completion
- [ ] Streaming responses
- [ ] Context cancellation and timeout
- [ ] Unified error handling

### Phase 2 — Provider Ecosystem

- [ ] OpenAI
- [ ] Anthropic
- [ ] Google Gemini
- [ ] DeepSeek
- [ ] Ollama
- [ ] Additional providers

### Phase 3 — Advanced Features

- [ ] Tool Calling
- [ ] Embeddings
- [ ] Model aliases
- [ ] Automatic retries
- [ ] Fallback strategies
- [ ] Load balancing
- [ ] Model routing

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
