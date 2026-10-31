# SDK

Module path: `github.com/icosmos-space/llmatrix/sdk`

Target: cover the LiteLLM **Python SDK** surface (not the Proxy gateway).

## Functions

| LiteLLM | LLMatrix |
|---|---|
| `completion` | `sdk.Completion` / `Client.ChatCompletion` |
| streaming completion | `sdk.CompletionStream` |
| `embedding` | `sdk.Embedding` |
| `image_generation` | `sdk.ImageGeneration` |
| `speech` | `sdk.Speech` |
| `transcription` | `sdk.Transcription` |
| `rerank` | `sdk.Rerank` |
| `Router` | `sdk.NewRouter` |
| callbacks | `WithSuccessCallback` / `WithFailureCallback` |
| exception mapping | `MapError` → `APIError` |

```go
sdk.Configure(sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
resp, err := sdk.Completion(ctx, &chat.Request{
    Model: "gpt-4o-mini", // alias → openai/gpt-4o-mini
    Messages: []chat.Message{{Role: "user", Content: "Hello"}},
})
```
