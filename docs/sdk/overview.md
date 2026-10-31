# SDK

Module path: `github.com/icosmos-space/llmatrix/sdk`

Inspired by LiteLLM’s Python SDK: one client API, provider selected from the model string.

## Phase 1

- `ChatCompletion` / `ChatCompletionStream` / `Embeddings`
- Model refs: `openai/...`, `deepseek/...`, `ollama/...`
- Shared adapters live in `engine/providers` (reused by Proxy later)

```go
client := sdk.NewClient(
    sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
    sdk.WithModel("openai/gpt-4o-mini"),
)
resp, err := client.ChatCompletion(ctx, &chat.Request{
    Messages: []chat.Message{{Role: "user", Content: "Hello"}},
})
```
