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
- LiteLLM SDK surface expansion: package-level `Completion`/`Embedding`/`ImageGeneration`/`Speech`/`Transcription`/`Rerank`,
  model aliases, callbacks, OpenAI-style `APIError`, Router load balancing, richer chat params,
  and OpenAI-compatible image/speech/transcription/rerank providers.
- OpenAI Responses API (`sdk.Responses`) with chat-bridge for non-native providers.
- Responses streaming (`ResponsesStream`) plus get/delete/cancel store helpers.
- Azure OpenAI provider (`azure/<deployment>`).
- Gemini function/tool calling.
- Groq / Together / Mistral OpenAI-compatible providers.
- Prompt caching: Anthropic `cache_control` content blocks + OpenAI `prompt_cache_key`.
- LiteLLM SDK parity: `Moderation`, `ImageEdit`, `TokenCounter` / `CompletionCost`,
  chat `reasoning_effort` / `max_completion_tokens` / `logprobs` / Anthropic `thinking`,
  Responses `previous_response_id`.
- In-memory response `Cache` (`WithMemoryCache`, request `Caching=true`) for chat/embedding/text.
- `TextCompletion` (OpenAI `/completions`).
- Files + Batches: `CreateFile` / `GetFile` / `DeleteFile` / `ListFiles` / `GetFileContent`,
  `CreateBatch` / `RetrieveBatch` / `CancelBatch` / `ListBatches`.
- OpenRouter / xAI / Fireworks OpenAI-compatible providers.
- Disk cache (`WithDiskCache`), `TextCompletionStream`.
- LiteLLM helpers: `SupportsFunctionCalling` / `SupportsVision` / `SupportsParallelFunctionCalling`,
  `GetSupportedOpenAIParams`, `GetModelInfo`.
- Client `WithOrganization`, `WithExtraHeaders`, `WithDropParams`.

