# LLMatrix SDK

Unified Go client inspired by [LiteLLM](https://github.com/BerriAI/litellm): one API, multiple providers.

Supported providers:

| Provider | Model example | Env key |
|---|---|---|
| OpenAI | `openai/gpt-4o-mini` | `OPENAI_API_KEY` |
| DeepSeek | `deepseek/deepseek-chat` | `DEEPSEEK_API_KEY` |
| Ollama | `ollama/llama3` | (optional) |
| Anthropic | `anthropic/claude-sonnet-4-5` | `ANTHROPIC_API_KEY` |
| Gemini | `gemini/gemini-2.0-flash` | `GEMINI_API_KEY` / `GOOGLE_API_KEY` |
| Generic | `openai-compatible/...` + `WithBaseURL` | `OPENAI_API_KEY` |

Non-streaming calls retry on rate limits / transient errors by default (`WithMaxRetries`, `WithRetry`).

Fallback models (tried after primary fails):

```go
client := sdk.NewClient(
    sdk.WithModel("openai/gpt-4o"),
    sdk.WithFallback("deepseek/deepseek-chat", "ollama/llama3"),
)
```

Tool calling:

```go
resp, err := client.ChatCompletion(ctx, &chat.Request{
    Messages: msgs,
    Tools:    []tools.Tool{weatherTool},
    ToolChoice: "auto",
})
```

## Install

```bash
go get github.com/icosmos-space/llmatrix/sdk
```

## Quick start

```go
client := sdk.NewClient(
    sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
    sdk.WithModel("openai/gpt-4o-mini"),
)

resp, err := client.ChatCompletion(ctx, &chat.Request{
    Messages: []chat.Message{
        {Role: "user", Content: "Hello!"},
    },
})
```

Streaming:

```go
stream, err := client.ChatCompletionStream(ctx, &chat.Request{Messages: msgs})
for stream.Next(ctx) {
    if ev := stream.Event(); ev.Type == streaming.EventTypeDelta {
        fmt.Print(ev.Delta)
    }
}
```

## Architecture note

SDK depends on `engine` for provider adapters so Proxy can reuse the same implementations later.
