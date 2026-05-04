package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/zadewu/focus/internal/domain"
)

var _ domain.FocusRepository = (*Repository)(nil)

type Repository struct {
	dir string
}

func NewRepository(dir string) *Repository {
	return &Repository{dir: dir}
}

func (r *Repository) run(args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = r.dir
	out, err := c.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func (r *Repository) Init() error {
	_, err := r.run("init", r.dir)
	if err != nil {
		return fmt.Errorf("init repository: %w", err)
	}
	return nil
}

func (r *Repository) Current() (string, error) {
	out, err := r.run("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get current focus: %w", err)
	}
	return out, nil
}

func (r *Repository) Exists(name string) bool {
	_, err := r.run("rev-parse", "--verify", name)
	return err == nil
}

func (r *Repository) List() ([]domain.Focus, error) {
	out, err := r.run("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, fmt.Errorf("list focuses: %w", err)
	}
	if out == "" {
		return nil, nil
	}

	lines := strings.Split(out, "\n")
	focuses := make([]domain.Focus, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		f := domain.Focus{}
		if rest, ok := strings.CutPrefix(line, "archive/"); ok {
			f.Name = rest
			f.Archived = true
		} else {
			f.Name = line
			f.Archived = false
		}
		f.CreatedAt = r.firstCommitTime(line)
		focuses = append(focuses, f)
	}
	return focuses, nil
}

// firstCommitTime returns the timestamp of the first commit on branch, zero time on error.
func (r *Repository) firstCommitTime(branch string) time.Time {
	c := exec.Command("git", "log", "--reverse", "--format=%ci", branch)
	c.Dir = r.dir
	out, err := c.Output()
	if err != nil {
		return time.Time{}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", lines[0])
	if err != nil {
		return time.Time{}
	}
	return t
}

func (r *Repository) Create(name string) error {
	// If no branches exist yet, use --orphan to avoid needing a parent commit.
	out, _ := r.run("branch", "--format=%(refname:short)")
	hasBranches := strings.TrimSpace(out) != ""

	if !hasBranches {
		if _, err := r.run("checkout", "--orphan", name); err != nil {
			return fmt.Errorf("create focus %q: %w", name, err)
		}
		if _, err := r.run("commit", "--allow-empty", "-m", "init"); err != nil {
			return fmt.Errorf("create focus %q: initial commit: %w", name, err)
		}
	} else {
		if _, err := r.run("checkout", "-b", name); err != nil {
			return fmt.Errorf("create focus %q: %w", name, err)
		}
	}
	return nil
}

func (r *Repository) Switch(name string) error {
	if _, err := r.run("checkout", name); err != nil {
		return fmt.Errorf("switch to focus %q: %w", name, err)
	}
	return nil
}

func (r *Repository) Archive(name string) error {
	if _, err := r.run("branch", "-m", name, "archive/"+name); err != nil {
		return fmt.Errorf("archive focus %q: %w", name, err)
	}
	return nil
}

func (r *Repository) AddNote(msg string) error {
	if _, err := r.run("commit", "--allow-empty", "-m", msg); err != nil {
		return fmt.Errorf("add note: %w", err)
	}
	return nil
}

func (r *Repository) RemoteGet(name string) (string, error) {
	out, err := r.run("remote", "get-url", name)
	if err != nil {
		return "", fmt.Errorf("get remote %q: %w", name, err)
	}
	return out, nil
}

func (r *Repository) RemoteSet(name, url string) error {
	if _, setErr := r.run("remote", "set-url", name, url); setErr != nil {
		if _, addErr := r.run("remote", "add", name, url); addErr != nil {
			return fmt.Errorf("set remote %q (set-url: %v): %w", name, setErr, addErr)
		}
	}
	return nil
}

func (r *Repository) PushAll(remote string) error {
	if _, err := r.run("push", "--all", remote); err != nil {
		return fmt.Errorf("push to %q: %w", remote, err)
	}
	if _, err := r.run("push", "--tags", remote); err != nil {
		return fmt.Errorf("push tags to %q: %w", remote, err)
	}
	return nil
}

func (r *Repository) FetchAll(remote string) error {
	if _, err := r.run("fetch", remote, "--prune"); err != nil {
		return fmt.Errorf("fetch from %q: %w", remote, err)
	}
	return nil
}

func (r *Repository) CheckoutRemoteBranches(remote string) error {
	out, err := r.run("branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return fmt.Errorf("list remote branches: %w", err)
	}
	prefix := remote + "/"
	for _, ref := range strings.Split(out, "\n") {
		ref = strings.TrimSpace(ref)
		if ref == "" || strings.HasSuffix(ref, "/HEAD") || !strings.HasPrefix(ref, prefix) {
			continue
		}
		branch := strings.TrimPrefix(ref, prefix)
		if r.Exists(branch) {
			continue
		}
		if _, err := r.run("checkout", "-b", branch, "--track", ref); err != nil {
			return fmt.Errorf("checkout %q: %w", branch, err)
		}
	}
	return nil
}

// CreateBranch creates a branch at current HEAD.
// For non-empty repos, the branch is created without switching.
// For empty repos, falls back to orphan checkout (HEAD changes).
func (r *Repository) CreateBranch(name string) error {
	out, _ := r.run("branch", "--format=%(refname:short)")
	if strings.TrimSpace(out) == "" {
		if _, err := r.run("checkout", "--orphan", name); err != nil {
			return fmt.Errorf("create branch %q: %w", name, err)
		}
		if _, err := r.run("commit", "--allow-empty", "-m", "imported"); err != nil {
			return fmt.Errorf("create branch %q: initial commit: %w", name, err)
		}
		return nil
	}
	if _, err := r.run("branch", name); err != nil {
		return fmt.Errorf("create branch %q: %w", name, err)
	}
	return nil
}

// RenameBranch renames a branch (works even if it is the currently checked-out branch).
func (r *Repository) RenameBranch(oldName, newName string) error {
	if _, err := r.run("branch", "-m", oldName, newName); err != nil {
		return fmt.Errorf("rename branch %q → %q: %w", oldName, newName, err)
	}
	return nil
}

func (r *Repository) GetNotes(name string) ([]domain.Note, error) {
	out, err := r.run("log", name, "--first-parent", "--pretty=format:%ci|%s")
	if err != nil {
		return nil, fmt.Errorf("get notes for %q: %w", name, err)
	}
	if out == "" {
		return nil, nil
	}

	lines := strings.Split(out, "\n")
	notes := make([]domain.Note, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tsStr, msg, ok := strings.Cut(line, "|")
		if !ok {
			continue
		}
		t, err := time.Parse("2006-01-02 15:04:05 -0700", tsStr)
		if err != nil {
			t = time.Time{}
		}
		notes = append(notes, domain.Note{Timestamp: t, Message: msg})
	}
	return notes, nil
}
