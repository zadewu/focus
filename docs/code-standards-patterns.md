# Focus: Code Patterns Reference

**Part of:** [code-standards.md](code-standards.md)

---

## Git Operations Pattern

Always exec the `git` binary via a wrapper — never use `go-git`:

```go
// internal/git/git.go
func Run(dir string, args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    out, err := cmd.CombinedOutput()
    return strings.TrimSpace(string(out)), err
}

func CurrentBranch(dir string) (string, error) {
    return Run(dir, "symbolic-ref", "--short", "HEAD")
}

func BranchExists(dir, name string) bool {
    _, err := Run(dir, "rev-parse", "--verify", name)
    return err == nil
}
```

Rationale: git binary is always present (repo needs it), simpler than go-git, no API churn.

---

## Terminal UI Pattern

Centralize all lipgloss styles in `internal/ui/ui.go`. Never define styles inline in cmd files:

```go
package ui

import "github.com/charmbracelet/lipgloss"

var (
    ActiveStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))  // green
    ArchivedStyle  = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("8"))  // gray
    CurrentMark    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))  // cyan
    HeaderStyle    = lipgloss.NewStyle().Bold(true).Underline(true)
    TimestampStyle = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("8"))
    NoteStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)
```

One style per semantic role. Cmd handlers call `ui.PrintFocusList()`, `ui.PrintLog()` — never format directly.

---

## $EDITOR Pattern

```go
func openEditor(initial string) (string, error) {
    editor := os.Getenv("EDITOR")
    if editor == "" { editor = "vi" }

    f, err := os.CreateTemp("", "focus-note-*.txt")
    if err != nil { return "", err }
    defer os.Remove(f.Name())
    f.WriteString(initial)
    f.Close()

    cmd := exec.Command(editor, f.Name())
    cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
    if err := cmd.Run(); err != nil { return "", err }

    data, err := os.ReadFile(f.Name())
    return string(data), err
}
```

---

## ZIP Workspace Files Pattern

```go
func zipDir(srcDir, outZip string) error {
    zf, err := os.Create(outZip)
    if err != nil { return err }
    defer zf.Close()
    w := zip.NewWriter(zf)
    defer w.Close()

    return filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
        if err != nil || info.IsDir() { return err }
        rel, _ := filepath.Rel(srcDir, path)
        zh, _ := zip.FileInfoHeader(info)
        zh.Name = filepath.ToSlash(rel)
        zh.Method = zip.Deflate
        fw, _ := w.CreateHeader(zh)
        f, _ := os.Open(path)
        defer f.Close()
        _, err = io.Copy(fw, f)
        return err
    })
}
```

---

## Testing Conventions

- Test files: `filename_test.go` in same package
- Function names: `TestFunctionName(t *testing.T)`
- Subtests: `t.Run("case description", func(t *testing.T) { ... })`
- No mocking git — test against a real temp git repo created in `t.TempDir()`

```go
func TestCreateBranch(t *testing.T) {
    dir := t.TempDir()
    exec.Command("git", "init", dir).Run()
    // ... test against real git repo
}
```

---

## Comment Standards

Exported functions: one-line doc comment:
```go
// CreateBranch creates a new git branch in the focus repo.
func CreateBranch(dir, name string) error { ... }
```

Unexported helpers: brief inline comment only when non-obvious:
```go
// expandHome replaces leading ~ with the user's home directory.
func expandHome(path string) string { ... }
```

No multi-paragraph docstrings. No comments restating what the code does.

---

## Pre-commit Checklist

- [ ] `go fmt ./...`
- [ ] `go vet ./...`
- [ ] `go mod tidy`
- [ ] `go test ./...` — all pass
- [ ] `go build ./...` — clean
- [ ] No file exceeds 200 LOC
- [ ] Errors include context (`fmt.Errorf("op: %w", err)`)
- [ ] No secrets or credentials committed
