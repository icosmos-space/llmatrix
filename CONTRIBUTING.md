# Contributing to LLMatrix

Thanks for your interest in contributing. Please open an issue before large architectural changes.

## Development setup

Requirements:

- Go 1.23+
- [Task](https://taskfile.dev/) (optional) or `make`

Clone and sync the workspace:

```bash
git clone https://github.com/icosmos-space/llmatrix.git
cd llmatrix
go work sync
```

## Modules

| Module | Path | Purpose |
|---|---|---|
| SDK | `sdk/` | Public client API |
| Engine | `engine/` | Provider adapters, routing, execution |
| Proxy | `proxy/` | HTTP proxy service and billing |

## Common commands

```bash
# tidy all modules
make tidy
# or: task tidy

# unit tests
make test

# lint
make lint

# build proxy binary
make build-proxy
```

## Pull request checklist

1. Keep changes scoped to one concern when possible.
2. Add or update tests for behavioral changes.
3. Run `make tidy`, `make test`, and `make lint` before opening a PR.
4. Update docs under `docs/` when public APIs or architecture change.
5. Follow conventional Go style and keep public APIs documented.

## Commit style

Prefer concise, imperative subject lines focused on *why*:

```text
sdk: add streaming event types for chat completions
engine: register deepseek as openai-compatible provider
```

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
