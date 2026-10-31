# SDK

Module path: `github.com/icosmos-space/llmatrix/sdk`

Inspired by LiteLLM’s Python SDK: one client API, provider selected from the model string.

## Providers

| Model prefix | Backend |
|---|---|
| `openai/` | OpenAI Chat Completions |
| `deepseek/` | DeepSeek (OpenAI-compatible) |
| `ollama/` | Ollama (OpenAI-compatible) |
| `anthropic/` | Anthropic Messages API |
| `gemini/` | Google Gemini generateContent |
| `openai-compatible/` | Custom OpenAI-compatible base URL |

```go
client := sdk.NewClient(
    sdk.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
    sdk.WithModel("anthropic/claude-sonnet-4-5"),
    sdk.WithMaxRetries(3),
)
resp, err := client.ChatCompletion(ctx, &chat.Request{
    Messages: []chat.Message{{Role: "user", Content: "Hello"}},
})
```
