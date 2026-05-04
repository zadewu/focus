# Focus: Project Roadmap

**Status:** Complete (MVP shipped)  
**Release Date:** 2026-05-03  
**Module:** `github.com/zadewu/focus`

---

## Phases

| # | Phase | Status | Completed | Complexity |
|---|-------|--------|-----------|-----------|
| 1 | Project Scaffold + Auto-Init | Complete | 2026-05-03 | Low |
| 2 | Core Focus Management | Complete | 2026-05-03 | Medium |
| 3 | Notes & Log | Complete | 2026-05-03 | Low |
| 4 | Status, Workspace & Config | Complete | 2026-05-03 | Low |
| 5 | Export: Markdown & Obsidian | Complete | 2026-05-03 | High |
| 6 | Polish, Build & Docs | Complete | 2026-05-03 | Low |
| 7 | Backup & Migration (Remote Push/Pull) | Complete | 2026-05-04 | Medium |
| 8 | Search across sessions | Complete | 2026-05-04 | Medium |

**All phases completed:** Go CLI fully functional with all planned features + backup/migration support + cross-session search.

---

## Phase Summaries

**Phase 1** — Bootstrap Go module, cobra root, `ensureInit()` for `~/.focus/`, `internal/git` + `internal/config` packages. ✓

**Phase 2** — `focus new`, `switch`, `list`, `archive`. Adds `internal/workspace`. Timestamp-based branch naming: `YYYY-MM-DD-HHmm__name`. ✓

**Phase 3** — `focus note` (inline + $EDITOR), `focus log`. Core data primitive: empty commits as notes. Short-name resolution for all commands. ✓

**Phase 4** — lipgloss styling in `internal/ui`, root status display, `focus workspace`, `focus config`. ✓

**Phase 5** — `focus export` (markdown) + `focus export --obsidian` (vault file using full branch name + journal append + non-md zip). ✓

**Phase 6** — `focus shell-init` for bash/zsh/fish/pwsh. Error handling pass, `go vet` clean, docs complete, working tree stays empty verification. ✓

**Phase 7** — `focus remote [url]` (get/set origin), `focus push` (backup to remote), `focus pull [--restore]` (fetch + restore branches for migration). ✓

**Phase 8** — `focus search <keyword>` (cross-session note search with rg → grep → Go fallback). ✓

---

## Risks

| Risk | Mitigation |
|------|-----------|
| git invocation platform-specific | Test macOS + Linux in Phase 2 |
| Obsidian vault structure varies | Support custom journal pattern via config |
| Branch name validation too strict/loose | Allow alphanumeric + `-_` only |
| `$EDITOR` not set | Fall back to `vi`, clear error on failure |

---

## Timeline

```
Week 1: Phase 1 (Mon-Tue) + Phase 2 (Wed-Fri)
Week 2: Phase 3 (Mon-Wed) + Phase 4 (Thu-Fri)
Week 3: Phase 5 — markdown export + obsidian core
Week 4: Phase 5 — journal + zip + Phase 6 (Thu-Fri)

Release: 2026-05-31
```

---

## Success Metrics

- All 10 commands (including `shell-init`) work end-to-end
- `go test ./...` passes, >80% coverage
- Short-name resolution works for all commands (no ambiguity handling errors in normal case)
- Operations complete <500ms
- `~/.focus/` working tree stays empty throughout
- `go install github.com/zadewu/focus@latest` works
- Shell integration auto-cd works for bash, zsh, fish, and pwsh
