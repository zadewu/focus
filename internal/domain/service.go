package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FocusService struct {
	repo      FocusRepository
	config    ConfigStore
	workspace WorkspaceStore
}

func NewFocusService(repo FocusRepository, config ConfigStore, workspace WorkspaceStore) *FocusService {
	return &FocusService{repo: repo, config: config, workspace: workspace}
}

func (s *FocusService) WorkspaceRoot() (string, error) {
	root, err := s.config.Get("workspace-root")
	if err != nil || root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "focus-workspaces"), nil
	}
	return root, nil
}

// NewFocus creates a new focus branch and provisions its workspace directory.
// The branch and workspace are named YYYY-MM-DD-HHmm__{name}.
// Returns the full branch name and the workspace path.
func (s *FocusService) NewFocus(name string) (fullName, wsPath string, err error) {
	if err := ValidateName(name); err != nil {
		return "", "", err
	}
	fullName = GenerateFullName(name, time.Now())
	if s.repo.Exists(fullName) {
		return "", "", fmt.Errorf("focus %q already exists", name)
	}
	if err := s.repo.Create(fullName); err != nil {
		return "", "", fmt.Errorf("create focus: %w", err)
	}
	if s.workspace != nil {
		wsPath, err = s.workspace.Ensure(fullName)
		return fullName, wsPath, err
	}
	return fullName, "", nil
}

// resolveFullName returns the full branch name for a user-supplied name.
// Accepts both full names (2026-05-03-2125__my-task) and short names (my-task).
// Returns an error if the name is not found or is ambiguous.
func (s *FocusService) resolveFullName(name string) (string, error) {
	if s.repo.Exists(name) {
		return name, nil
	}
	focuses, err := s.repo.List()
	if err != nil {
		return "", err
	}
	var matches []string
	for _, f := range focuses {
		if ExtractShortName(f.Name) == name {
			if f.Archived {
				matches = append(matches, "archive/"+f.Name)
			} else {
				matches = append(matches, f.Name)
			}
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("focus %q not found", name)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("ambiguous focus name %q — matches: %s", name, strings.Join(matches, ", "))
	}
}

// SwitchFocus checks out an existing focus branch and returns its workspace path.
// Accepts both full names and short names. Archived focuses cannot be switched to.
func (s *FocusService) SwitchFocus(name string) (wsPath string, err error) {
	fullName, err := s.resolveFullName(name)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(fullName, "archive/") {
		return "", fmt.Errorf("focus %q is archived; cannot switch to an archived focus", name)
	}
	if err := s.repo.Switch(fullName); err != nil {
		return "", fmt.Errorf("switch focus: %w", err)
	}
	if s.workspace != nil {
		wsName, _ := strings.CutPrefix(fullName, "archive/")
		return s.workspace.Path(wsName), nil
	}
	return "", nil
}

// ListFocuses returns all focus sessions and the name of the currently active one.
func (s *FocusService) ListFocuses() (focuses []Focus, current string, err error) {
	focuses, err = s.repo.List()
	if err != nil {
		return nil, "", err
	}
	current, _ = s.repo.Current()
	return focuses, current, nil
}

// ArchiveFocus renames a focus branch to archive/<name>.
// If name is empty, the currently active focus is archived.
// When archiving the active branch, it switches to another active focus first.
func (s *FocusService) ArchiveFocus(name string) error {
	if name == "" {
		current, err := s.repo.Current()
		if err != nil {
			return fmt.Errorf("no active focus")
		}
		name = current
	} else {
		var err error
		name, err = s.resolveFullName(name)
		if err != nil {
			return err
		}
		if bare, isArchived := strings.CutPrefix(name, "archive/"); isArchived {
			return fmt.Errorf("focus %q is already archived", ExtractShortName(bare))
		}
	}
	if !s.repo.Exists(name) {
		return fmt.Errorf("focus %q not found", name)
	}

	// If archiving the currently active branch, switch away first.
	current, _ := s.repo.Current()
	if name == current {
		focuses, _ := s.repo.List()
		for _, f := range focuses {
			if !f.Archived && f.Name != name {
				_ = s.repo.Switch(f.Name)
				break
			}
		}
		// If no other active focus exists, proceed anyway;
		// Archive renames the branch even in detached HEAD state.
	}

	return s.repo.Archive(name)
}

// AddNote appends a timestamped note to the current active focus.
func (s *FocusService) AddNote(msg string) error {
	current, err := s.repo.Current()
	if err != nil || current == "" {
		return fmt.Errorf("no active focus — run: focus new <name>")
	}
	if strings.HasPrefix(current, "archive/") {
		return fmt.Errorf("current focus is archived; switch to an active focus first")
	}
	return s.repo.AddNote(msg)
}

// Status returns the current focus name, workspace path, and recent notes for status display.
func (s *FocusService) Status() (current, wsPath string, notes []Note, err error) {
	current, err = s.repo.Current()
	if err != nil || current == "" {
		return "", "", nil, nil // no active focus — not an error
	}
	if s.workspace != nil {
		wsPath = s.workspace.Path(current)
	}
	notes, _ = s.repo.GetNotes(current)
	return current, wsPath, notes, nil
}

// WorkspacePath returns the workspace directory path for a named focus (defaults to current).
// Accepts both full names and short names.
func (s *FocusService) WorkspacePath(name string) (string, error) {
	if name == "" {
		var err error
		name, err = s.repo.Current()
		if err != nil || name == "" {
			return "", fmt.Errorf("no active focus")
		}
	} else {
		var err error
		name, err = s.resolveFullName(name)
		if err != nil {
			return "", err
		}
	}
	if s.workspace == nil {
		return "", fmt.Errorf("workspace not configured")
	}
	wsName, _ := strings.CutPrefix(name, "archive/")
	return s.workspace.Path(wsName), nil
}

// GetConfig retrieves a configuration value by key.
func (s *FocusService) GetConfig(key string) (string, error) {
	v, err := s.config.Get(key)
	if err != nil || v == "" {
		return "", fmt.Errorf("key not set: %s", key)
	}
	return v, nil
}

// SetConfig stores a configuration value by key.
func (s *FocusService) SetConfig(key, value string) error {
	return s.config.Set(key, value)
}

// Export gathers focus data and delegates to the provided exporter.
func (s *FocusService) Export(name string, exporter Exporter) error {
	focus, notes, files, err := s.resolveExportData(name)
	if err != nil {
		return err
	}
	return exporter.Export(focus, notes, files)
}

func (s *FocusService) resolveExportData(name string) (Focus, []Note, []File, error) {
	if name == "" {
		current, err := s.repo.Current()
		if err != nil || current == "" {
			return Focus{}, nil, nil, fmt.Errorf("no active focus")
		}
		name = current
	} else {
		var err error
		name, err = s.resolveFullName(name)
		if err != nil {
			return Focus{}, nil, nil, err
		}
	}
	displayName, _ := strings.CutPrefix(name, "archive/")
	focuses, _ := s.repo.List()
	var focus Focus
	for _, f := range focuses {
		if f.Name == displayName {
			focus = f
			break
		}
	}
	if focus.Name == "" {
		focus = Focus{Name: displayName}
	}
	notes, err := s.repo.GetNotes(name)
	if err != nil {
		return Focus{}, nil, nil, err
	}
	var files []File
	if s.workspace != nil {
		files, _ = s.workspace.ListFiles(displayName)
	}
	return focus, notes, files, nil
}

// GetLog returns the display name and notes for a given focus (defaults to current).
// Accepts both full names and short names.
func (s *FocusService) GetLog(name string) (focusName string, notes []Note, err error) {
	if name == "" {
		name, err = s.repo.Current()
		if err != nil || name == "" {
			return "", nil, fmt.Errorf("no active focus — run: focus new <name>")
		}
	} else {
		name, err = s.resolveFullName(name)
		if err != nil {
			return "", nil, err
		}
	}
	displayName, _ := strings.CutPrefix(name, "archive/")
	notes, err = s.repo.GetNotes(name)
	return displayName, notes, err
}

const defaultRemote = "origin"

func (s *FocusService) RemoteGet() (string, error) {
	url, err := s.repo.RemoteGet(defaultRemote)
	if err != nil {
		return "", fmt.Errorf("no remote configured — run: focus remote <url>")
	}
	return url, nil
}

func (s *FocusService) RemoteSet(url string) error {
	return s.repo.RemoteSet(defaultRemote, url)
}

// Push pushes all branches and tags to the configured remote.
// Returns the remote URL so callers can display it without an extra lookup.
func (s *FocusService) Push() (string, error) {
	url, err := s.repo.RemoteGet(defaultRemote)
	if err != nil {
		return "", fmt.Errorf("no remote configured — run: focus remote <url>")
	}
	return url, s.repo.PushAll(defaultRemote)
}

// Pull fetches from remote. If restore is true, local tracking branches are created
// for every remote branch not already present locally (migration path).
func (s *FocusService) Pull(restore bool) error {
	if _, err := s.repo.RemoteGet(defaultRemote); err != nil {
		return fmt.Errorf("no remote configured — run: focus remote <url>")
	}
	if err := s.repo.FetchAll(defaultRemote); err != nil {
		return err
	}
	if restore {
		return s.repo.CheckoutRemoteBranches(defaultRemote)
	}
	return nil
}
