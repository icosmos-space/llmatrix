# LLMatrix SDK

Unified Go client inspired by [LiteLLM](https://github.com/BerriAI/litellm): one API, multiple providers.

## LiteLLM SDK parity

| LiteLLM | LLMatrix |
|---|---|
| `completion()` | `sdk.Completion` / `Client.ChatCompletion` |
| `completion(..., stream=True)` | `sdk.CompletionStream` / `Client.ChatCompletionStream` |
| `embedding()` | `sdk.Embedding` / `Client.Embeddings` |
| `image_generation()` | `sdk.ImageGeneration` / `Client.ImageGeneration` |
| `speech()` | `sdk.Speech` / `Client.Speech` |
| `transcription()` | `sdk.Transcription` / `Client.Transcription` |
| `rerank()` | `sdk.Rerank` / `Client.Rerank` |
| `Router` | `sdk.NewRouter` |
| retries / fallbacks | `WithRetry` / `WithFallback` |
| model aliases | built-in (`gpt-4o-mini` → `openai/gpt-4o-mini`) |
| success/failure callbacks | `WithSuccessCallback` / `WithFailureCallback` |
| exception mapping | `sdk.MapError` / `APIError` |
| `responses()` | Planned |

## Providers

| Provider | Model example | Env key |
|---|---|---|
| OpenAI | `openai/gpt-4o-mini` or alias `gpt-4o-mini` | `OPENAI_API_KEY` |
| DeepSeek | `deepseek/deepseek-chat` | `DEEPSEEK_API_KEY` |
| Ollama | `ollama/llama3` | (optional) |
| Anthropic | `anthropic/claude-sonnet-4-5` | `ANTHROPIC_API_KEY` |
| Gemini | `gemini/gemini-2.0-flash` | `GEMINI_API_KEY` / `GOOGLE_API_KEY` |
| Generic | `openai-compatible/...` + `WithBaseURL` | `OPENAI_API_KEY` |

## Quick start

```go
import "github.com/icosmos-space/llmatrix/sdk"
import "github.com/icosmos-space/llmatrix/sdk/chat"

sdk.Configure(sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

resp, err := sdk.Completion(ctx, &chat.Request{
    Model: "gpt-4o-mini",
    Messages: []chat.Message{{Role: "user", Content: "Hello"}},
})
```

Client style:

```go
client := sdk.NewClient(
    sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
    sdk.WithModel("openai/gpt-4o-mini"),
    sdk.WithFallback("deepseek/deepseek-chat"),
)
```

Router (load balance deployments):

```go
r := sdk.NewRouter([]sdk.Deployment{
    {Name: "chat", Provider: "openai", Model: "gpt-4o-mini", Weight: 1},
    {Name: "chat", Provider: "openai", Model: "gpt-4o", Weight: 1},
})
resp, err := r.ChatCompletion(ctx, "chat", &chat.Request{Messages: msgs})
```
