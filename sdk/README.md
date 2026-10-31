# LLMatrix SDK

Public Go client for LLMatrix.

Use this module to call providers directly through the Engine, or to talk to a
deployed LLMatrix Proxy with the same request/response shapes.

## Install

```bash
go get github.com/icosmos-space/llmatrix/sdk
```

## Packages

| Package | Role |
|---|---|
| `sdk` | Client entry point, options, config, errors |
| `sdk/chat` | Chat completion helpers |
| `sdk/embedding` | Embedding helpers |
| `sdk/streaming` | Stream event types |
| `sdk/tools` | Tool / function-calling helpers |
| `sdk/proxy` | Client for LLMatrix Proxy |

## Examples

See `sdk/examples/` for direct, proxy, streaming, and tools samples.
