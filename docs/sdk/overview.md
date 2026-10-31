# SDK

Module path: `github.com/icosmos-space/llmatrix/sdk`

Target: cover the LiteLLM **Python SDK** surface (not the Proxy gateway).

## Functions

| LiteLLM | LLMatrix |
|---|---|
| `completion` | `sdk.Completion` / `Client.ChatCompletion` |
| streaming completion | `sdk.CompletionStream` |
| `text_completion` | `sdk.TextCompletion` |
| `embedding` | `sdk.Embedding` |
| `image_generation` | `sdk.ImageGeneration` |
| `image_edit` | `sdk.ImageEdit` |
| image variations | `sdk.ImageVariation` |
| `speech` | `sdk.Speech` |
| `transcription` | `sdk.Transcription` |
| `rerank` | `sdk.Rerank` |
| `moderation` | `sdk.Moderation` |
| `create_file` / file helpers | `sdk.CreateFile` / `GetFile` / `DeleteFile` / `ListFiles` / `GetFileContent` |
| `create_batch` / batch helpers | `sdk.CreateBatch` / `RetrieveBatch` / `CancelBatch` / `ListBatches` |
| `token_counter` | `sdk.TokenCounter` / `TokenCounterMessages` |
| `completion_cost` | `sdk.CompletionCost` / `CompletionCostFromUsage` |
| `responses` | `sdk.Responses` (OpenAI native; others bridged via chat) |
| streaming responses | `sdk.ResponsesStream` |
| get/delete/cancel responses | `sdk.GetResponse` / `DeleteResponse` / `CancelResponse` |
| `Cache()` | `WithMemoryCache` / `WithDiskCache` / `WithRedisCache` + request `Caching` |
| `litellm_proxy/` | `litellm_proxy/<model>` or `WithUseLiteLLMProxy(true)` |
| assistants / threads / runs | `CreateAssistants` / `CreateThread` / `AddMessage` / `CreateRun` / … |
| `supports_function_calling` | `sdk.SupportsFunctionCalling` |
| `supports_vision` | `sdk.SupportsVision` |
| `supports_parallel_function_calling` | `sdk.SupportsParallelFunctionCalling` |
| `get_supported_openai_params` | `sdk.GetSupportedOpenAIParams` |
| `get_model_info` | `sdk.GetModelInfo` |
| `drop_params` | `WithDropParams(true)` |
| `Router` | `sdk.NewRouter` |
| callbacks | `WithSuccessCallback` / `WithFailureCallback` |
| exception mapping | `MapError` → `APIError` |
| prompt caching | `chat.Message.CacheControl` / `Parts` (Anthropic); `PromptCacheKey` (OpenAI) |

```go
resp, err := sdk.Responses(ctx, &responses.Request{
    Model: "openai/gpt-4o-mini",
    Input: "Tell me a short story",
})
fmt.Println(resp.OutputText())
```

Azure:

```go
client := sdk.NewClient(
    sdk.WithProvider("azure"),
    sdk.WithAPIKey(os.Getenv("AZURE_API_KEY")),
    sdk.WithBaseURL(os.Getenv("AZURE_API_BASE")),
    sdk.WithAPIVersion("2024-08-01-preview"),
    sdk.WithModel("azure/my-gpt4o-deployment"),
)
```
