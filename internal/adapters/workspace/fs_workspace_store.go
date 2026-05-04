package workspace

import (
	"os"
	"path/filepath"

	"github.com/zadewu/focus/internal/domain"
)

// Compile-time check: FSWorkspaceStore must implement domain.WorkspaceStore.
var _ domain.WorkspaceStore = (*FSWorkspaceStore)(nil)

// FSWorkspaceStore manages focus workspace directories on the local filesystem.
type FSWorkspaceStore struct {
	root string
}

// NewFSWorkspaceStore creates a workspace store rooted at the given directory.
func NewFSWorkspaceStore(root string) *FSWorkspaceStore {
	return &FSWorkspaceStore{root: root}
}

// Path returns the absolute path to the workspace directory for the given focus name.
func (s *FSWorkspaceStore) Path(name string) string {
	return filepath.Join(s.root, name)
}

// Ensure creates the workspace directory for the given focus name if it does not exist.
func (s *FSWorkspaceStore) Ensure(name string) (string, error) {
	p := s.Path(name)
	return p, os.MkdirAll(p, 0o755)
}

// ListFiles returns all files in the workspace directory for the given focus name.
// Returns nil (no error) if the directory does not exist.
func (s *FSWorkspaceStore) ListFiles(name string) ([]domain.File, error) {
	dir := s.Path(name)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []domain.File
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			// Skip unreadable files; they may be locked or permission-denied.
			continue
		}
		files = append(files, domain.File{Name: e.Name(), Content: content})
	}
	return files, nil
}
