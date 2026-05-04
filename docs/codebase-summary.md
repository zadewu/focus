# Focus: Codebase Structure

**Status:** Implemented (MVP complete + UI/UX polish)  
**Architecture:** Hexagonal (ports & adapters)  
**Language:** Go 1.22+  
**Module:** `github.com/zadewu/focus`  
**Last updated:** 2026-05-04

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
│   ├── import.go                              # focus import [--dry-run]
│   ├── remote.go                              # focus remote [url]
│   ├── push.go                                # focus push
│   ├── pull.go                                # focus pull [--restore]
│   ├── search.go                              # focus search <keyword>
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
            ├── lipgloss_renderer.go           # Terminal renderer (lipgloss + adaptive colours)
            ├── interactive_list.go            # Bubble Tea TUI (scroll + fuzzy filter)
            └── terminal_utils.go              # word-wrap + terminal width detection
```

---

## Domain Layer (`internal/domain/`)

### `focus.go` — Entities
- `Focus` struct: `Name`, `CreatedAt`, `Archived`
- `Note` struct: `Timestamp`, `Message`
- `SearchResult` struct: `Focus`, `Note` — pairs a note with its session context
- `ValidateName()` — business rule: no spaces/slashes, not starting with `archive`

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
    CreateBranch(name string) error
    RenameBranch(oldName, newName string) error
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
- `SearchNotes(keyword)` — filter notes across all sessions using rg → grep → Go fallback
- `Export(name, exporter)` — collect notes + files, call exporter
- `RemoteGet()` — get configured origin URL
- `RemoteSet(url)` — set origin URL
- `Push()` — push all branches + tags to origin
- `Pull(restore)` — fetch from origin; create local tracking branches if restore=true
- `ImportFocuses(dryRun)` — migrate legacy focus sessions to canonical name format (two-pass: branches, then workspace dirs)

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
Styles for status, list, log output using adaptive colours (readable on both light and dark terminals). Not an interface-backed adapter — imported directly by `cmd/`.

### `ui/interactive_list.go`
Bubble Tea TUI model for `focus list` interactive mode (when run in a TTY). Provides scrollable, fuzzy-filterable list with `/` filter toggle, Enter to select/switch, `q`/Esc to cancel. Falls back to plain list when piped.

### `ui/terminal_utils.go`
Terminal utilities: `getTerminalWidth()` detects terminal column count, `wordWrap()` breaks long notes across lines with indentation aligned to message start. Used by `PrintLog()` and `PrintStatus()`.

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
- **cmd/**: 15 command handlers (new, switch, list, archive, note, log, workspace, config, export, import, search, shell-init, remote, push, pull)
  - `list.go` — interactive TUI when connected to TTY; falls back to plain list when piped
  - `log.go` — displays notes with word-wrapped output
- **domain/**: focus.go, service.go, ports.go fully implemented
  - `SearchResult` struct, `SearchNotes()` method — cross-session search with pluggable backends
  - `isCurrentPrefixed()` — detect canonical YYYY-MM-DD-HHmm__ prefix
  - `ParseImportName()` — convert legacy names to canonical format
  - `ImportFocuses()` — two-pass migration (branches, then workspace dirs)
- **adapters/**: git, config, workspace, export (markdown + obsidian), and ui all functional
  - `CreateBranch()`, `RenameBranch()` added to git adapter
  - **UI enhancements**: Adaptive colour palette (readable on light/dark/transparent terminals), interactive Bubble Tea TUI for `focus list`, word-wrap for long notes
  - `PrintSearchResults()` in ui — renders search matches with focus/note context
  - `RunInteractiveList()` — scrollable, fuzzy-filterable session selector
  - `wordWrap()` — intelligent line breaking with indentation alignment
- **Tests**: Comprehensive unit test coverage for domain logic + remote operations + import migration
- **Deployment**: Ready for `go install github.com/zadewu/focus@latest`
