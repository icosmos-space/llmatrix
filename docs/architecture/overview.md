# Architecture Overview

LLMatrix is organized as a Go workspace with three primary modules:

```text
Application / SDK Client
        |
        v
   +---------+
   |   SDK   |  public client API
   +---------+
        |
        +------------------+
        |                  |
        v                  v
   +---------+       +-----------+
   | Engine  | <---- |   Proxy   |  HTTP service
   +---------+       +-----------+
        |
   providers / protocol / router
```

- **SDK** — idiomatic Go client for apps.
- **Engine** — provider adapters, routing, and execution.
- **Proxy** — multi-tenant HTTP gateway with auth, billing, and channel management.
