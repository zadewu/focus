# Dependency Injection Rules

* **DI-1 (MUST)** Perform wiring only in `/cmd` or `/platform`
* **DI-2 (MUST)** Use constructor injection
* **DI-3 (MUST)** MUST NOT instantiate adapters inside core
* **DI-4 (SHOULD)** Keep wiring explicit (avoid hidden magic)

## Example

```go
repo := postgres.NewUserRepository(db)
uc := usecase.NewCreateUser(repo)
handler := http.NewHandler(uc)
```
