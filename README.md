# focus

A terminal tool for managing personal work sessions. All session metadata lives in git (`~/.focus/`) — zero files in the working tree. Real work lives in `~/focus-workspaces/<name>/`.

## Install

**Linux & macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/zadewu/focus/main/install.sh | sh
```

The script detects your OS and architecture, downloads the binary, verifies the SHA256 checksum, and installs to `~/.local/bin` (or `/usr/local/bin` if that's not writable). Set `FOCUS_VERSION=vX.Y.Z` to pin a specific release.

**Via Go:**

```bash
go install github.com/zadewu/focus@latest
```

Binaries and checksums for each release: [github.com/zadewu/focus/releases](https://github.com/zadewu/focus/releases)

## Quick Start

```bash
focus new my-task          # create session + workspace
focus note "first thought" # add a note
focus log                  # view notes
focus list                 # list all sessions
focus archive my-task      # archive when done
```

## Commands

| Command | Description |
|---------|-------------|
| `focus` | Show current session + recent notes |
| `focus new <name>` | Create session and workspace |
| `focus switch <name>` | Switch active session |
| `focus list` | List all sessions |
| `focus archive [name]` | Archive a session |
| `focus note [msg]` | Add a note (opens $EDITOR if no msg) |
| `focus log [name]` | Show note history |
| `focus workspace [name]` | Print workspace path |
| `focus config <key> [value]` | Get/set config |
| `focus export [name]` | Export to markdown |
| `focus export [name] --obsidian` | Export to Obsidian vault |
| `focus import [--dry-run]` | Migrate legacy sessions to canonical format |
| `focus remote [url]` | Get or set backup remote URL |
| `focus push` | Push all sessions to remote |
| `focus pull [--restore]` | Fetch sessions from remote |

## Obsidian Integration

```bash
focus config obsidian-vault ~/Documents/my-vault
focus export --obsidian        # writes vault file + journal entry + attachments zip
```

## Backup & Migration

```bash
# Set a remote (GitHub or any Git host) and push all sessions
focus remote https://github.com/you/focus-data.git
focus push

# On a new machine: clone the repo, then restore local branches
git clone https://github.com/you/focus-data.git ~/.focus/
focus pull --restore
```

`focus pull` without `--restore` fetches remote state (updates remote-tracking refs).
`focus pull --restore` also creates a local branch for every remote branch — use this after migrating to a new machine.

## Migration: Upgrading from Legacy Sessions

If you have old sessions created with earlier versions of `focus`, their names may use legacy formats (e.g., `YYYY-MM-DD--name` or `plain-name`). The new canonical format is `YYYY-MM-DD-HHmm__name`.

**Migrate all sessions:**

```bash
# Preview changes without modifying anything
focus import --dry-run

# Apply the migration (renames branches + workspace directories)
focus import
```

The `import` command runs two passes:
1. Renames legacy git branches in `~/.focus/`
2. Renames legacy workspace directories in `~/focus-workspaces/` and creates missing branches

Name conversion:
- `YYYY-MM-DD--my-task` → `YYYY-MM-DD-0000__my-task` (HHmm defaults to 0000)
- `plain-name` → `2000-01-01-0000__plain-name` (plain names get sentinel date)
- `YYYY-MM-DD-HHmm__my-task` → unchanged (already canonical)

After migration, `focus list` will display all sessions, and `focus switch my-task` will work with both full and short names.

## Shell Integration

Source once to auto-`cd` into the workspace on `focus new` / `focus switch`:

```bash
# bash / zsh — add to ~/.bashrc or ~/.zshrc
eval "$(focus shell-init)"

# fish — add to ~/.config/fish/config.fish
focus shell-init | source

# PowerShell — add to $PROFILE
Invoke-Expression (focus shell-init --shell pwsh)
```

To preview the script without sourcing: `focus shell-init --shell fish`

## How It Works

- **Sessions** = branches in `~/.focus/` (a git repo with an empty working tree)
- **Notes** = empty git commits with the message as the note body
- **Config** = `git config focus.*` in `~/.focus/.git/config`
- **Workspaces** = plain directories in `~/focus-workspaces/<name>/`

All metadata is queryable with standard git tools: `git -C ~/.focus log`, `git -C ~/.focus branch -a`.
