# Focus: System Architecture

**Status:** Design (pre-implementation)  
**Last updated:** 2026-05-03

---

## High-Level Architecture

Focus is a two-layer system: **metadata layer** (git) + **workspace layer** (filesystem).

```
┌──────────────────────────────────────────��──────────────────────┐
│  User Shell (bash/zsh/fish)                                     │
│  $ focus new debug-auth                                         │
│  $ focus note "found the bug"                                   │
│  $ fcd debug-auth  (cd ~/focus-workspaces/debug-auth)          │
└──────────┬──────────────────────────────────────────────────────┘
           │
           ↓
┌─────────────────────────────────────────────────────────────────┐
│  Focus CLI (Go binary)                                          │
│  cmd/: root, new, switch, list, archive, note, log,            │
│        workspace, config, export                                │
│  internal/: git, config, workspace, export, ui                  │
└──────────┬──────────────────┬──────────────────────────────────┘
           │                  │
           ↓                  ↓
    ┌────────────┐    ┌─────────────────────┐
    │ ~/.focus/  │    │ ~/focus-workspaces/ │
    │ (git repo) │    │ debug-auth/         │
    │ refs/heads │    │   notes.md          │
    │   debug-   │    │   repro.sh          │
    │   archive/ │    │ planning-v2/        │
    └────────────┘    └─────────────────────┘
```

---

## Core Data Flows

### 1. Create Focus Session
```
focus new debug-auth
  → ValidateName("debug-auth")
  → git checkout -b debug-auth  (in ~/.focus/)
  → mkdir ~/focus-workspaces/debug-auth
  → print: "Workspace: ~/focus-workspaces/debug-auth"
```

### 2. Add Note
```
focus note "found the bug"
  → check HEAD branch exists + is not archived
  → git commit --allow-empty -m "found the bug"
```

### 3. Export to Obsidian
```
focus export --obsidian
  → read git log on current branch (commits = notes)
  → read workspace .md files
  → write <vault>/Focus/YYYY-MM-DD-HHmm__<name>.md
  → zip non-.md workspace files → <vault>/Focus/attachments/<name>.zip
  → if <vault>/01 Daily/YYYY/MM/YYYY-MM-DD.md exists:
      append/update ## Focus section with [[wikilinks]]
    else: warn + skip
```

---

## Hexagonal Architecture

```
         ┌──────────────────────────────────┐
         │         cmd/ (Primary Adapter)   │
         │     cobra commands → FocusService│
         └──────────────┬───────────────────┘
                        │ calls
                        ↓
         ┌──────────────────────────────────┐
         │      internal/domain/            │
         │  FocusService  (use cases)        │
         │  Focus, Note   (entities)         │
         │  Ports         (interfaces)       │
         └──────────────┬───────────────────┘
                        │ implemented by
                        ↓
         ┌──────────────────────────────────┐
         │    internal/adapters/ (Secondary)│
         │  git/     → FocusRepository      │
         │  config/  → ConfigStore           │
         │  workspace/ → WorkspaceStore     │
         │  export/  → Exporter             │
         │  ui/      → terminal renderer    │
         └──────────────────────────────────┘
```

**Dependency rule:** `cmd` → `domain` ← `adapters`. Domain has zero imports from adapters or cmd.

## Package Responsibilities

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Primary adapter | `cmd/*` | Parse args, call FocusService, print output |
| Domain | `domain/service.go` | Use cases; pure business logic |
| Domain | `domain/ports.go` | Port interfaces (FocusRepository, ConfigStore, etc.) |
| Domain | `domain/focus.go` | Focus entity, Note value object, name validation |
| Secondary adapter | `adapters/git/` | FocusRepository via `exec git` |
| Secondary adapter | `adapters/config/` | ConfigStore via `git config focus.*` |
| Secondary adapter | `adapters/workspace/` | WorkspaceStore via filesystem |
| Secondary adapter | `adapters/export/` | Exporter implementations (markdown, obsidian) |
| Secondary adapter | `adapters/ui/` | lipgloss terminal renderer |

See [system-architecture-storage.md](system-architecture-storage.md) for storage model, config keys, and Obsidian vault structure.
