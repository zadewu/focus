# Focus: Project Roadmap

**Status:** Planning phase  
**Target Release:** 2026-05-31  
**Module:** `github.com/zadewu/focus`

---

## Phases

| # | Phase | Status | Est. Days | Complexity | Depends On |
|---|-------|--------|-----------|-----------|-----------|
| 1 | Project Scaffold + Auto-Init | Pending | 2 | Low | — |
| 2 | Core Focus Management | Pending | 5 | Medium | 1 |
| 3 | Notes & Log | Pending | 3 | Low | 1 |
| 4 | Status, Workspace & Config | Pending | 3 | Low | 2, 3 |
| 5 | Export: Markdown & Obsidian | Pending | 8 | High | 3, 4 |
| 6 | Polish, Build & Docs | Pending | 2 | Low | 5 |

**Total:** ~23 days (4 weeks)

Plan files: `plans/260503-1211-focus-terminal-app-go-cli/`

---

## Phase Summaries

**Phase 1** — Bootstrap Go module, cobra root, `ensureInit()` for `~/.focus/`, `internal/git` + `internal/config` packages.

**Phase 2** — `focus new`, `switch`, `list`, `archive`. Adds `internal/workspace`. Branch name validation.

**Phase 3** — `focus note` (inline + $EDITOR), `focus log`. Core data primitive: empty commits as notes.

**Phase 4** — lipgloss styling in `internal/ui`, root status display, `focus workspace`, `focus config`.

**Phase 5** — `focus export` (markdown) + `focus export --obsidian` (vault file + journal append + non-md zip).

**Phase 6** — Error handling pass, `go vet` clean, README, working tree stays empty verification.

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

- All 10 commands work end-to-end
- `go test ./...` passes, >80% coverage
- Operations complete <500ms
- `~/.focus/` working tree stays empty throughout
- `go install github.com/zadewu/focus@latest` works
