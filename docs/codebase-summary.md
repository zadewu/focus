# Focus: Codebase Structure

**Status:** Implemented (MVP complete)  
**Architecture:** Hexagonal (ports & adapters)  
**Language:** Go 1.22+  
**Module:** `github.com/zadewu/focus`  
**Last updated:** 2026-05-03

---

## Directory Tree

```
focus/
├── main.go                                    # Binary entrypoint
├── go.mod                                     # github.com/zadewu/focus
├── go.sum
│
├── cmd/                                       # PRIMARY ADAPTER — CLI (cobra)
│   ├── root.go                                # focus (status display + auto-init)
│   ├── new.go                                 # focus new <name>
│   ├── switch.go                              # focus switch <name>
│   ├── list.go                                # focus list
│   ├── archive.go                             # focus archive [name]
│   ├── note.go                                # focus note [msg]
│   ├── log.go                                 # focus log [name]
│   ├── workspace.go                           # focus workspace [name]
│   ├── config.go                              # focus config <key> <value>
│   ├── export.go                              # focus export [--obsidian]
│   ├── remote.go                              # focus remote [url]
│   ├── push.go                                # focus push
│   ├── pull.go                                # focus pull [--restore]
│   └── shell_init.go                          # focus shell-init
│
└── internal/
    ├── domain/                                # DOMAIN CORE — no external deps
    │   ├── focus.go                           # Focus entity, Note value object
    │   ├── service.go                         # Use cases (FocusService)
    │   └── ports.go                           # Port interfaces
    │
    └── adapters/                              # SECONDARY ADAPTERS
        ├── git/
        │   └── git_repository.go              # FocusRepository via git exec
        ├── config/
        │   └── git_config_store.go            # ConfigStore via git config
        ├── workspace/
        │   └── fs_workspace_store.go          # WorkspaceStore via filesystem
        ├── export/
        │   ├── markdown_exporter.go           # Exporter → plain markdown
        │   └── obsidian_exporter.go           # Exporter → Obsidian vault
        └── ui/
            └── lipgloss_renderer.go           # Terminal renderer (lipgloss)
```

---

## Domain Layer (`internal/domain/`)

### `focus.go` — Entities
- `Focus` struct: `Name`, `CreatedAt`, `IsArchived`
- `Note` struct: `Timestamp`, `Message`
- `validateName()` — business rule: no spaces/slashes, not starting with `archive`

### `ports.go` — Port Interfaces
```go
type FocusRepository interface {
    Init() error
    Create(name string) error
    Switch(name string) error
    Archive(name string) error
    List() ([]Focus, error)
    Current() (string, error)
    AddNote(msg string) error
    GetNotes(name string) ([]Note, error)
    Exists(name string) bool
    RemoteGet(name string) (string, error)
    RemoteSet(name, url string) error
    PushAll(remote string) error
    FetchAll(remote string) error
    CheckoutRemoteBranches(remote string) error
}

type ConfigStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

type WorkspaceStore interface {
    Path(name string) string
    Ensure(name string) (string, error)
    ListFiles(name string) ([]File, error)
}

type Exporter interface {
    Export(focus Focus, notes []Note, files []File) error
}
```

### `service.go` — Use Cases (FocusService)
Accepts port interfaces via constructor injection:
- `CreateFocus(name)` — validate + repo.Create + workspace.Ensure
- `SwitchFocus(name)` — validate exists + repo.Switch
- `AddNote(msg)` — validate not archived + repo.AddNote
- `ArchiveFocus(name)` — repo.Archive
- `ListFocuses()` — repo.List
- `GetNotes(name)` — repo.GetNotes
- `Export(name, exporter)` — collect notes + files, call exporter
- `RemoteGet()` — get configured origin URL
- `RemoteSet(url)` — set origin URL
- `Push()` — push all branches + tags to origin
- `Pull(restore)` — fetch from origin; create local tracking branches if restore=true

---

## Adapters Layer (`internal/adapters/`)

### `git/git_repository.go` — `FocusRepository`
Implements domain port using `exec.Command("git", ...)`. Sets `cmd.Dir = ~/.focus/`.

### `config/git_config_store.go` — `ConfigStore`
Reads/writes `git config focus.*` in `~/.focus/`. Provides typed getters with defaults in service layer (not adapter).

### `workspace/fs_workspace_store.go` — `WorkspaceStore`
Creates/reads `~/focus-workspaces/<name>/`. Lists files split by `.md` / non-`.md`.

### `export/markdown_exporter.go` — `Exporter`
Renders notes + workspace `.md` files to a markdown file in CWD.

### `export/obsidian_exporter.go` — `Exporter`
Writes `<vault>/Focus/YYYY-MM-DD-HHmm__<name>.md`, zips non-`.md` files, appends journal.

### `ui/lipgloss_renderer.go`
Styles for status, list, log output. Not an interface-backed adapter — imported directly by `cmd/`.

---

## Wire-Up (`main.go` / `cmd/root.go`)

```go
repo      := git.NewRepository(focusDir)
cfg       := config.NewGitConfigStore(focusDir)
ws        := workspace.NewFSStore(cfg)
service   := domain.NewFocusService(repo, cfg, ws)
// inject service into cmd handlers
```

---

## Key Patterns

- **Dependency direction:** `cmd` → `domain` ← `adapters` (adapters depend on domain ports, never vice versa)
- **No circular deps:** domain has zero imports from adapters or cmd
- **Errors:** domain returns domain errors; adapters wrap git/fs errors with context
- **Config defaults:** defined in `FocusService`, not in adapters

---

## Implementation Status

All source files implemented and tested:
- **cmd/**: 13 command handlers (new, switch, list, archive, note, log, workspace, config, export, shell-init, remote, push, pull)
- **domain/**: focus.go, service.go, ports.go fully implemented
- **adapters/**: git, config, workspace, export (markdown + obsidian), and ui all functional
- **Tests**: Comprehensive unit test coverage for domain logic + remote operations
- **Deployment**: Ready for `go install github.com/zadewu/focus@latest`
