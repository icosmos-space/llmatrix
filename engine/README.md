# LLMatrix Engine

Core execution layer: unified types, provider interfaces, registry, routing,
protocol adapters, and concrete provider implementations.

## Layout

| Package | Role |
|---|---|
| `types` | Shared request / response models |
| `provider` | Provider abstractions |
| `registry` | Provider registration and factories |
| `engine` | Call execution |
| `router` | Model routing strategies |
| `protocol` | Wire-format converters (OpenAI / Anthropic / Gemini) |
| `providers` | Concrete provider clients |
| `model` / `usage` / `stream` / `errors` | Supporting primitives |

SDK and Proxy both depend on Engine for provider execution.
