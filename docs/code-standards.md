# Focus: Code Standards

**Language:** Go 1.22+  
**Last updated:** 2026-05-03

---

## File Organization

- **Go files:** `snake_case.go` (e.g. `git.go`, `workspace.go`, `markdown.go`)
- **Max file size:** 200 LOC — split into sibling modules if exceeded
- **Package names:** lowercase, no underscores (`git`, `config`, `ui`)
- **One responsibility per package**

```
cmd/                          → Primary adapter: cobra command handlers
internal/domain/              → Domain core: entities, use cases, port interfaces
internal/adapters/git/        → Secondary adapter: FocusRepository via git exec
internal/adapters/config/     → Secondary adapter: ConfigStore via git config
internal/adapters/workspace/  → Secondary adapter: WorkspaceStore via filesystem
internal/adapters/export/     → Secondary adapter: Exporter (markdown, obsidian)
internal/adapters/ui/         → Terminal renderer (lipgloss)
```

**Dependency rule (mandatory):** `cmd` → `domain` ← `adapters`. The domain package must never import from `cmd/` or `adapters/`.

---

## Naming Conventions

| Scope | Convention | Example |
|-------|-----------|---------|
| Exported functions | PascalCase | `CreateBranch()`, `WorkspaceRoot()` |
| Unexported functions | camelCase | `sanitizeName()`, `expandHome()` |
| Local variables | camelCase | `focusName`, `branchList` |
| Constants | PascalCase | `DefaultWorkspaceRoot` |

---

## Error Handling

Always wrap errors with context using `%w`:

```go
func CreateBranch(dir, name string) error {
    if err := validateName(name); err != nil {
        return fmt.Errorf("create branch: %w", err)
    }
    _, err := git.Run(dir, "checkout", "-b", name)
    return err
}
```

User-facing errors must suggest remediation:
```go
return fmt.Errorf("obsidian vault not configured\nRun: focus config obsidian-vault <path>")
```

---

## Cobra Command Pattern

```go
var newCmd = &cobra.Command{
    Use:   "new <name>",
    Short: "Create a new focus session",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := git.CreateBranch(focusDir, args[0]); err != nil {
            return err
        }
        ws, err := workspace.Ensure(focusDir, args[0])
        if err != nil { return err }
        fmt.Printf("Focus: %s\nWorkspace: %s\n", args[0], ws)
        return nil
    },
}
```

Rules:
- Always use `RunE` (not `Run`) — cobra handles printing returned errors
- Validate arg count with `cobra.ExactArgs(N)` / `cobra.MaximumNArgs(N)`
- Return errors; don't `os.Exit()` from commands

---

## Input Validation

```go
func validateName(name string) error {
    if name == "" { return errors.New("name cannot be empty") }
    if strings.ContainsAny(name, " \t\n/\\") {
        return fmt.Errorf("invalid characters in name: %q", name)
    }
    if strings.HasPrefix(name, "archive") {
        return errors.New("name cannot start with 'archive'")
    }
    if len(name) > 100 { return errors.New("name too long (max 100 chars)") }
    return nil
}
```

---

## Approved Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling |
| stdlib only | Everything else |

**Forbidden:** `go-git` (use exec git instead), ORM libraries, unnecessary HTTP clients.

See [code-standards-patterns.md](code-standards-patterns.md) for git ops, UI patterns, testing, and pre-commit checklist.
