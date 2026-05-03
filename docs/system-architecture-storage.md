# Focus: Storage Model & Obsidian Integration

**Part of:** [system-architecture.md](system-architecture.md)

---

## Git Repository (~/.focus/)

```
~/.focus/
├── .git/
│   ├── refs/heads/
│   │   ├── debug-auth           ← active focus (branch)
│   │   ├── planning-v2
│   │   └── archive/old-task     ← archived focus
│   └── config                   ← focus.* config keys
└── (empty working tree, no files ever checked out)
```

| Concept | Git primitive |
|---------|--------------|
| Focus name | Branch name |
| Note | Empty commit message |
| Note timestamp | Commit date (automatic) |
| Active focus | HEAD |
| Archived focus | Branch prefixed `archive/` |
| App config | `git config focus.*` |

**Key git commands used:**
- Create: `git checkout -b <name>`
- Switch: `git checkout <name>`
- Add note: `git commit --allow-empty -m "<msg>"`
- Archive: `git branch -m <name> archive/<name>`
- Read notes: `git log <branch> --pretty=format:%ci|%s`

---

## Workspace Directories

```
~/focus-workspaces/
├── debug-auth/
│   ├── notes.md          ← user-created markdown (exported to vault)
│   ├── repro.sh          ← non-md (zipped on obsidian export)
│   └── cloned-repo/      ← arbitrary structure, user owns it
└── planning-v2/
    └── v2-design.md
```

- App creates directory only; never manages contents
- Preserved on `focus archive` (not deleted)
- Exported as: `.md` files → appended to vault focus file; non-`.md` → zipped

---

## Config Storage

**Location:** `~/.focus/.git/config`

```ini
[focus]
    workspace-root = /Users/huy/focus-workspaces
    obsidian-vault = /Users/huy/Obsidian/Personal
    obsidian-journal-pattern = 01 Daily/{YYYY}/{MM}/{YYYY}-{MM}-{DD}
```

| Key | Default |
|-----|---------|
| `focus.workspace-root` | `~/focus-workspaces` |
| `focus.obsidian-vault` | (none — required for export) |
| `focus.obsidian-journal-pattern` | `01 Daily/{YYYY}/{MM}/{YYYY}-{MM}-{DD}` |

---

## Obsidian Vault Structure

```
<vault>/
├── Focus/
│   ├── 2026-05-03-0900__debug-auth.md    ← focus file (creation timestamp)
│   └── attachments/
│       └── debug-auth.zip                ← non-md workspace files
└── 01 Daily/2026/05/2026-05-03.md        ← user's daily journal
```

**Focus file content** (`Focus/YYYY-MM-DD-HHmm__<name>.md`):
```markdown
# Focus: debug-auth

**Created:** 2026-05-03 09:00

## Notes

- **09:00** — Started investigation
- **14:22** — Found token validation bug

## Workspace Files

### notes.md
[contents]
```

**Daily journal append** (if file exists — never created by focus):
```markdown
## Focus

- [[Focus/2026-05-03-0900__debug-auth|debug-auth]] — 2 notes today
  - 09:30 — Started investigation
  - 14:22 — Found token validation bug
```

---

## Error Handling Strategy

| Level | Example | Response |
|-------|---------|----------|
| Git op fails | branch already exists | `"create branch: already exists"` |
| Config missing | obsidian-vault not set | `"obsidian vault not configured; run: focus config obsidian-vault <path>"` |
| File system | can't create workspace | `"create workspace: permission denied"` |
| Journal missing | daily file not found | warn + skip (not an error) |

Errors propagate: internal packages → cmd → cobra → stderr + exit 1.
