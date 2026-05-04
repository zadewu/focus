# Focus: Project Overview & Product Development Requirements

## What is Focus?

**Focus** is a personal terminal application for managing work sessions. It solves context-switching overhead by treating each session as a first-class entity with its own metadata (creation date, notes) and workspace directory. Built in Go, it runs as a single binary with no external dependencies.

## Problem Statement

Knowledge workers juggle multiple parallel tasks (debugging, planning, feature work, etc.). Current tools force choice: use git for everything (pollutes real repos) or use external apps (context lag). Focus fills the gap: git internals track *what/when*, filesystem holds *actual work*.

## User Personas

1. **Backend Developer** — Switches between 3-5 bug investigations. Needs quick context restoration and note history per bug.
2. **Researcher/Designer** — Long-running exploratory work. Wants daily journal integration and artifact archival.
3. **DevOps Engineer** — Tracks deployment tasks, runbooks, config changes. Exports to Obsidian for runbook sharing.

## Key Features

| Feature | Purpose |
|---------|---------|
| `focus new <name>` | Create session with auto-initialized workspace dir |
| `focus note "text"` | Log session event (creates empty commit with message) |
| `focus switch <name>` | Restore context by switching active session |
| `focus list` | View all active and archived sessions |
| `focus archive <name>` | Retire session (workspace preserved for later retrieval) |
| `focus export --obsidian` | Push notes + workspace files to Obsidian vault + daily journal |
| `focus shell-init` | Generate shell integration script for auto-cd on new/switch |
| Config persistence | Store vault path, workspace root via git config |

## Non-Goals

- **NOT a task manager** (no due dates, priorities, subtasks)
- **NOT a note app** (metadata is session labels + commit messages only)
- **NOT a file sync tool** (workspaces are local-only)
- **NOT a git client** (doesn't replace git for real repos)
- **NOT cloud-backed** (single-machine state)

## Success Criteria

1. Single `go install` command works
2. `focus new X` + workspace dir creation completes in <100ms
3. `focus export --obsidian` updates vault focus file + daily journal section in <1s
4. `focus list` output distinguishes active vs archived sessions
5. Zero files in `~/.focus/` working tree (only git metadata)
6. Shell `fcd <name>` integrates seamlessly with `cd`

## Target Users

- Solo developers, researchers, technical writers
- Anyone using Obsidian as personal knowledge base
- Teams wanting session-level metadata without git pollution

## Implementation Status

**MVP shipped:** 2026-05-03
- All 10 commands fully implemented and tested
- Hexagonal architecture with zero domain dependencies
- Git-backed session metadata via local `~/.focus/` repository
- Timestamped branch names (`YYYY-MM-DD-HHmm__name`) + workspace directories
- Short-name resolution for all commands (switch/archive/workspace/log accept both full and short names)
- Shell integration (`focus shell-init`) for auto-cd on new/switch
- Markdown and Obsidian vault export with journal integration
- Ready for production use: `go install github.com/zadewu/focus@latest`

## Dependencies

- **Go 1.22+** (language runtime)
- **github.com/spf13/cobra** (CLI framework)
- **github.com/charmbracelet/lipgloss** (terminal styling)
- **git** (system binary, assumed present)
- **$EDITOR** (optional, for note composition)

---

**Version:** 1.0 (implementation complete)  
**Last updated:** 2026-05-03  
**Module:** `github.com/zadewu/focus`
