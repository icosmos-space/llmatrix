# LLMatrix SDK

Unified Go client inspired by [LiteLLM](https://github.com/BerriAI/litellm): one API, multiple providers.

## LiteLLM SDK parity

| LiteLLM | LLMatrix |
|---|---|
| `completion()` | `sdk.Completion` / `Client.ChatCompletion` |
| `completion(..., stream=True)` | `sdk.CompletionStream` / `Client.ChatCompletionStream` |
| `text_completion()` | `sdk.TextCompletion` / `Client.TextCompletion` |
| `embedding()` | `sdk.Embedding` / `Client.Embeddings` |
| `image_generation()` | `sdk.ImageGeneration` / `Client.ImageGeneration` |
| `image_edit()` | `sdk.ImageEdit` / `Client.ImageEdit` |
| `speech()` | `sdk.Speech` / `Client.Speech` |
| `transcription()` | `sdk.Transcription` / `Client.Transcription` |
| `rerank()` | `sdk.Rerank` / `Client.Rerank` |
| `moderation()` | `sdk.Moderation` / `Client.Moderation` |
| `create_file()` / files | `CreateFile` / `GetFile` / `DeleteFile` / `ListFiles` / `GetFileContent` |
| `create_batch()` / batches | `CreateBatch` / `RetrieveBatch` / `CancelBatch` / `ListBatches` |
| `token_counter` | `sdk.TokenCounter` / `TokenCounterMessages` |
| `completion_cost` | `sdk.CompletionCost` / `CompletionCostFromUsage` |
| `Cache()` + `caching=True` | `WithMemoryCache` / `WithDiskCache` / `WithRedisCache` + request `Caching` |
| `litellm_proxy/` | `litellm_proxy/<model>` or `WithUseLiteLLMProxy(true)` |
| assistants / threads / runs | `CreateAssistants` / `CreateThread` / `AddMessage` / `CreateRun` / … |
| image variations | `sdk.ImageVariation` |
| `text_completion(..., stream=True)` | `sdk.TextCompletionStream` |
| `supports_function_calling` / `supports_vision` | `SupportsFunctionCalling` / `SupportsVision` |
| `get_supported_openai_params` | `GetSupportedOpenAIParams` |
| `get_model_info` | `GetModelInfo` |
| `drop_params` | `WithDropParams(true)` |
| `responses()` | `sdk.Responses` / `Client.Responses` (native OpenAI; else chat bridge) |
| streaming responses | `sdk.ResponsesStream` / `Client.ResponsesStream` |
| get/delete/cancel responses | `GetResponse` / `DeleteResponse` / `CancelResponse` |
| `Router` | `sdk.NewRouter` |
| retries / fallbacks | `WithRetry` / `WithFallback` |
| model aliases | built-in (`gpt-4o-mini` → `openai/gpt-4o-mini`) |
| success/failure callbacks | `WithSuccessCallback` / `WithFailureCallback` |
| exception mapping | `sdk.MapError` / `APIError` |

## Providers

| Provider | Model example | Env key |
|---|---|---|
| OpenAI | `openai/gpt-4o-mini` or alias `gpt-4o-mini` | `OPENAI_API_KEY` |
| DeepSeek | `deepseek/deepseek-chat` | `DEEPSEEK_API_KEY` |
| Ollama | `ollama/llama3` | (optional) |
| Anthropic | `anthropic/claude-sonnet-4-5` | `ANTHROPIC_API_KEY` |
| Gemini | `gemini/gemini-2.0-flash` | `GEMINI_API_KEY` / `GOOGLE_API_KEY` |
| Azure OpenAI | `azure/<deployment>` + `WithBaseURL` / `WithAPIVersion` | `AZURE_API_KEY` |
| Groq | `groq/llama-3.3-70b-versatile` | `GROQ_API_KEY` |
| Together | `together/<model-id>` | `TOGETHER_API_KEY` |
| Mistral | `mistral/mistral-large-latest` | `MISTRAL_API_KEY` |
| OpenRouter | `openrouter/<model>` | `OPENROUTER_API_KEY` |
| xAI | `xai/grok-2` | `XAI_API_KEY` |
| Fireworks | `fireworks/<model>` | `FIREWORKS_API_KEY` |
| Perplexity | `perplexity/<model>` | `PERPLEXITYAI_API_KEY` |
| Cohere | `cohere/<model>` | `COHERE_API_KEY` |
| LiteLLM Proxy | `litellm_proxy/<model>` | `LITELLM_PROXY_API_KEY` + `LITELLM_PROXY_API_BASE` |
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
