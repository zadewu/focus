# Adapter Rules

## Inbound Adapters (HTTP, gRPC, CLI)

* **AD-IN-1 (MUST)** Convert external request → usecase input
* **AD-IN-2 (MUST)** Call usecase via ports only
* **AD-IN-3 (MUST)** MUST NOT contain business logic
* **AD-IN-4 (MUST)** Validate input at boundary

## Outbound Adapters (DB, APIs)

* **AD-OUT-1 (MUST)** Implement ports defined in core
* **AD-OUT-2 (MUST)** Handle I/O concerns (timeouts, retries)
* **AD-OUT-3 (MUST)** Map infra models ↔ domain models
* **AD-OUT-4 (MUST)** MUST NOT leak DB/API models into core

## General

* **AD-GEN-1 (MUST)** Declare interface compliance:

```go
var _ UserRepository = (*PostgresUserRepository)(nil)
```

* **AD-GEN-2 (SHOULD)** Keep adapters thin and focused
* **AD-GEN-3 (SHOULD)** Avoid shared logic between adapters unless abstracted properly
