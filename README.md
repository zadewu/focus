# Focus

A personal terminal app to manage work sessions. Each session (focus) gets its own metadata (notes, timestamp) and workspace directory—all backed by git internals, exported to Obsidian on demand.

---

## Quick Install

```bash
go install github.com/zadewu/focus@latest
```

Requires Go 1.22+ and git.

---

## Quick Start

```bash
# Create a new focus session
focus new debug-auth-bug
cd ~/focus-workspaces/debug-auth-bug

# Add a note
focus note "Found the issue in token validation"

# Switch to another focus
focus switch planning-v2

# See all sessions
focus list

# Export to Obsidian
focus export --obsidian
```

---

## Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `focus` | Show current focus + recent notes | `focus` |
| `focus new <name>` | Create session + workspace dir | `focus new debug-auth` |
| `focus switch <name>` | Switch active session | `focus switch planning-v2` |
| `focus list` | List all sessions (active & archived) | `focus list` |
| `focus archive [name]` | Archive session (workspace preserved) | `focus archive debug-auth` |
| `focus note [msg]` | Add note (or open $EDITOR) | `focus note "found it"` |
| `focus log [name]` | Show notes for current or named session | `focus log` |
| `focus workspace [name]` | Print workspace directory path | `focus workspace` |
| `focus export [name]` | Export to markdown | `focus export` |
| `focus export --obsidian` | Export to Obsidian vault | `focus export --obsidian` |
| `focus config <key> <val>` | Set config value | `focus config obsidian-vault ~/Obsidian` |

---

## Configuration

Set config values with `focus config <key> <value>`:

| Key | Default | Purpose |
|-----|---------|---------|
| `focus.workspace-root` | `~/focus-workspaces` | Root directory for workspace dirs |
| `focus.obsidian-vault` | *(required for export)* | Path to Obsidian vault |
| `focus.obsidian-journal-pattern` | `01 Daily/{YYYY}/{MM}/{YYYY}-{MM}-{DD}` | Journal file path pattern |

**Example:**
```bash
focus config obsidian-vault ~/Obsidian/Personal
```

---

## How It Works

Focus uses a two-part design:

**Metadata layer** (`~/.focus/`) — A git repository with empty working tree
- Branches = session names
- Commits = notes (empty commits with message body)
- Config = settings (obsidian-vault, workspace-root, etc.)

**Workspace layer** (`~/focus-workspaces/<name>/`) — Regular directories
- User manages content (code, documents, notes, etc.)
- Preserved on archive
- Exported to Obsidian as markdown + attachments

When you run `focus note "text"`, it creates a git commit with your message on the current branch. When you export to Obsidian, it reads all commits (notes), renders them chronologically, and appends workspace markdown files to a vault focus file.

---

## Obsidian Integration

### Setup

```bash
focus config obsidian-vault /path/to/your/Obsidian/vault
```

### Export

```bash
focus export --obsidian
```

Creates:
- `<vault>/Focus/YYYY-MM-DD-HHmm__<name>.md` — Focus notes + workspace .md files
- `<vault>/Focus/attachments/<name>.zip` — Non-markdown workspace files
- Updates `<vault>/01 Daily/YYYY/MM/YYYY-MM-DD.md` — Appends ## Focus section with today's notes (if file exists)

### Example Journal Section

```markdown
## Focus

- [[Focus/2026-05-03-1430__debug-auth|debug-auth]] — 2 notes today
  - 14:35 — Found the token validation bug
  - 16:22 — Fix in place, testing now
- [[Focus/2026-05-03-1100__planning-v2|planning-v2]] — 1 note today
  - 11:00 — Drafted new API structure
```

---

## Shell Integration

Add this function to your `.bashrc`, `.zshrc`, or `config.fish`:

### Bash/Zsh
```bash
function fcd() { cd "$(focus workspace ${1:-.})"; }
```

### Fish
```fish
function fcd; cd (focus workspace $argv); end
```

**Usage:**
```bash
fcd debug-auth
# Same as: cd ~/focus-workspaces/debug-auth
```

---

## Documentation

- **[Project Overview](./docs/project-overview-pdr.md)** — Goals, user personas, feature list
- **[System Architecture](./docs/system-architecture.md)** — Two-part design, data flow, config storage
- **[Code Standards](./docs/code-standards.md)** — Go conventions, error handling, cobra patterns
- **[Codebase Structure](./docs/codebase-summary.md)** — Planned packages and modules
- **[Project Roadmap](./docs/project-roadmap.md)** — 6 phases, timeline, deliverables

---

## Status

**Current:** Design phase (documentation complete, implementation starting)

See [Project Roadmap](./docs/project-roadmap.md) for phased delivery plan (4 weeks estimated).

---

## License

[To be determined]

---

## Author

Designed and implemented by Huy Phan  
Project: `github.com/zadewu/focus`
