# Focus: System Architecture

**Status:** Implemented (MVP complete)  
**Last updated:** 2026-05-03

---

## High-Level Architecture

Focus is a two-layer system: **metadata layer** (git) + **workspace layer** (filesystem).

```
┌──────────────────────────────────────────��──────────────────────┐
│  User Shell (bash/zsh/fish)                                     │
│  $ focus new my-task                                            │
│  $ focus note "found the bug"                                   │
│  $ focus switch 2026-05-03-2125__my-task                        │
│  (auto-cd via: eval "$(focus shell-init)")                      │
└──────────┬──────────────────────────────────────────────────────┘
           │
           ↓
┌─────────────────────────────────────────────────────────────────┐
│  Focus CLI (Go binary)                                          │
│  cmd/: root, new, switch, list, archive, note, log,            │
│        workspace, config, export, shell-init                    │
│  internal/: git, config, workspace, export, ui                  │
└──────────┬──────────────────┬──────────────────────────────────┘
           │                  │
           ↓                  ↓
    ┌────────────┐    ┌─────────────────────┐
    │ ~/.focus/  │    │ ~/focus-workspaces/  │
    │ (git repo) │    │ 2026-05-03-2125__   │
    │ refs/heads │    │   my-task/          │
    │   2026-... │    │   notes.md          │
    │   archive/ │    │   repro.sh          │
    └────────────┘    └─────────────────────┘
```

---

## Core Data Flows

### 1. Create Focus Session
```
focus new my-task
  → ValidateName("my-task")
  → generate full name: 2026-05-03-2125__my-task
  → git checkout -b 2026-05-03-2125__my-task  (in ~/.focus/)
  → mkdir ~/focus-workspaces/2026-05-03-2125__my-task
  → print: "Created: 2026-05-03-2125__my-task"
           "Workspace: ~/focus-workspaces/2026-05-03-2125__my-task"
```

### 2. Switch Focus (accepts short or full names)
```
focus switch my-task
  OR
focus switch 2026-05-03-2125__my-task
  → resolveFullName("my-task") → "2026-05-03-2125__my-task"
  → git checkout 2026-05-03-2125__my-task
  → return workspace path for auto-cd
```

### 3. Add Note
```
focus note "found the bug"
  → check HEAD branch exists + is not archived
  → git commit --allow-empty -m "found the bug"
```

### 4. Export to Obsidian
```
focus export --obsidian
  → read git log on current branch (commits = notes)
  → read workspace .md files
  → write <vault>/Focus/2026-05-03-2125__my-task.md
  → zip non-.md workspace files → <vault>/Focus/attachments/2026-05-03-2125__my-task.zip
  → if <vault>/01 Daily/YYYY/MM/YYYY-MM-DD.md exists:
      append/update ## Focus section with [[wikilinks]]
    else: warn + skip
```

### 5. Shell Integration
```
eval "$(focus shell-init)"  # bash/zsh
focus shell-init | source   # fish

focus new my-task
  → creates branch & workspace
  → shell function auto-cd into workspace on success

focus switch 2026-05-03-2125__my-task
  → switches branch
  → shell function auto-cd into workspace on success
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
