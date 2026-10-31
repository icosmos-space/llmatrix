# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial Go workspace layout with `sdk`, `engine`, and `proxy` modules.
- Scaffolding for providers, protocol adapters, proxy service, docs, and CI workflows.
- Phase 1 LiteLLM-style SDK: model parsing (`provider/model`), chat / stream / embeddings
  via OpenAI-compatible HTTP for OpenAI, DeepSeek, and Ollama.
- Anthropic Messages API provider (chat + SSE stream).
- Gemini generateContent / streamGenerateContent / embedContent providers.
- `engine/retry` with exponential backoff; SDK applies it to chat + embeddings.
- Tool calling for OpenAI-compatible and Anthropic (`tools` / `tool_choice` / `tool_calls`).
- SDK `WithFallback` model chain; `engine/router` FallbackStrategy helpers.
- Proxy `ChatService` wired to Engine providers; SDK `proxy.Client` HTTP implementation.

