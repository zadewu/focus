# Error Handling Rules

## General

* **ERR-1 (MUST)** Wrap errors with context:

```go
fmt.Errorf("create user: %w", err)
```

* **ERR-2 (MUST)** Use `errors.Is` / `errors.As`
* **ERR-3 (MUST)** MUST NOT compare error strings

## Boundary Rules

* **ERR-4 (MUST)** Adapters MUST map infra errors → domain errors
* **ERR-5 (MUST)** Core MUST NOT depend on infra error types
* **ERR-6 (SHOULD)** Define domain-level errors clearly
