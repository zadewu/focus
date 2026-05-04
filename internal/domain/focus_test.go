package domain

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractShortName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"2026-05-03-2125__my-task", "my-task"},  // current format
		{"2026-05-03--my-task", "my-task"},        // legacy YYYY-MM-DD-- format
		{"some-task", "some-task"},                // plain name, no prefix
		{"2026-05-03-2125__", ""},                 // empty short name edge case
		{"2026-05-03--", ""},                      // empty legacy short name
		{"my--task", "my--task"},                   // "--" not at position 10
		{"feat-my-do--work", "feat-my-do--work"},   // non-digit prefix, no misclassification
	}
	for _, tc := range cases {
		got := ExtractShortName(tc.input)
		if got != tc.want {
			t.Errorf("ExtractShortName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// --- Service Layer Tests for Remote/Push/Pull ---

// MockRepository is a mock FocusRepository for testing service logic.
type MockRepository struct {
	remotes     map[string]string
	pushing     bool
	pulling     bool
	checkingOut bool
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		remotes: make(map[string]string),
	}
}

func (m *MockRepository) Init() error                                  { return nil }
func (m *MockRepository) Create(name string) error                    { return nil }
func (m *MockRepository) Switch(name string) error                    { return nil }
func (m *MockRepository) Archive(name string) error                   { return nil }
func (m *MockRepository) List() ([]Focus, error)                      { return nil, nil }
func (m *MockRepository) Current() (string, error)                    { return "main", nil }
func (m *MockRepository) AddNote(msg string) error                    { return nil }
func (m *MockRepository) GetNotes(name string) ([]Note, error)        { return nil, nil }
func (m *MockRepository) Exists(name string) bool                     { return true }
func (m *MockRepository) RemoteGet(name string) (string, error) {
	url, ok := m.remotes[name]
	if !ok {
		return "", fmt.Errorf("no such remote: %s", name)
	}
	return url, nil
}

func (m *MockRepository) RemoteSet(name, url string) error {
	m.remotes[name] = url
	return nil
}

func (m *MockRepository) PushAll(remote string) error {
	if _, ok := m.remotes[remote]; !ok {
		return fmt.Errorf("no such remote: %s", remote)
	}
	m.pushing = true
	return nil
}

func (m *MockRepository) FetchAll(remote string) error {
	if _, ok := m.remotes[remote]; !ok {
		return fmt.Errorf("no such remote: %s", remote)
	}
	m.pulling = true
	return nil
}

func (m *MockRepository) CheckoutRemoteBranches(remote string) error {
	if _, ok := m.remotes[remote]; !ok {
		return fmt.Errorf("no such remote: %s", remote)
	}
	m.checkingOut = true
	return nil
}

// MockConfigStore is a mock ConfigStore for testing.
type MockConfigStore struct {
	values map[string]string
}

func NewMockConfigStore() *MockConfigStore {
	return &MockConfigStore{values: make(map[string]string)}
}

func (m *MockConfigStore) Get(key string) (string, error) {
	v, ok := m.values[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return v, nil
}

func (m *MockConfigStore) Set(key, value string) error {
	m.values[key] = value
	return nil
}

// MockWorkspaceStore is a mock WorkspaceStore for testing.
type MockWorkspaceStore struct {
	paths map[string]string
}

func NewMockWorkspaceStore() *MockWorkspaceStore {
	return &MockWorkspaceStore{paths: make(map[string]string)}
}

func (m *MockWorkspaceStore) Path(name string) string {
	return m.paths[name]
}

func (m *MockWorkspaceStore) Ensure(name string) (string, error) {
	return "", nil
}

func (m *MockWorkspaceStore) ListFiles(name string) ([]File, error) {
	return nil, nil
}

// TestServicePush_NoRemote tests that Push() returns error when no remote configured.
func TestServicePush_NoRemote(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Push without setting remote should error
	_, err := svc.Push()
	if err == nil {
		t.Errorf("expected error when pushing without remote, got nil")
	}
	if !strings.Contains(err.Error(), "no remote configured") {
		t.Errorf("expected 'no remote configured' in error, got: %v", err)
	}
}

// TestServicePull_NoRemote tests that Pull() returns error when no remote configured.
func TestServicePull_NoRemote(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Pull without setting remote should error
	err := svc.Pull(false)
	if err == nil {
		t.Errorf("expected error when pulling without remote, got nil")
	}
	if !strings.Contains(err.Error(), "no remote configured") {
		t.Errorf("expected 'no remote configured' in error, got: %v", err)
	}
}

// TestServicePush_WithRemote tests that Push() succeeds when remote is configured.
func TestServicePush_WithRemote(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Set remote
	if err := repo.RemoteSet("origin", "https://example.com/repo.git"); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Push should succeed
	_, err := svc.Push()
	if err != nil {
		t.Errorf("Push with remote configured failed: %v", err)
	}

	// Verify push was called
	if !repo.pushing {
		t.Errorf("expected PushAll to be called")
	}
}

// TestServicePull_NoRestore tests Pull(false) fetches but does not checkout remote branches.
func TestServicePull_NoRestore(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Set remote
	if err := repo.RemoteSet("origin", "https://example.com/repo.git"); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Pull without restore
	err := svc.Pull(false)
	if err != nil {
		t.Errorf("Pull(false) failed: %v", err)
	}

	// Verify fetch was called
	if !repo.pulling {
		t.Errorf("expected FetchAll to be called")
	}
}

// TestServicePull_WithRestore tests Pull(true) fetches and checks out remote branches.
func TestServicePull_WithRestore(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Set remote
	if err := repo.RemoteSet("origin", "https://example.com/repo.git"); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Pull with restore
	err := svc.Pull(true)
	if err != nil {
		t.Errorf("Pull(true) failed: %v", err)
	}

	// Verify both fetch and checkout were called
	if !repo.pulling {
		t.Errorf("expected FetchAll to be called")
	}
	if !repo.checkingOut {
		t.Errorf("expected CheckoutRemoteBranches to be called when restore=true")
	}
}

// TestServiceRemoteSet_Success tests RemoteSet stores the remote URL.
func TestServiceRemoteSet_Success(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	testURL := "https://github.com/example/repo.git"
	if err := svc.RemoteSet(testURL); err != nil {
		t.Errorf("RemoteSet failed: %v", err)
	}

	// Verify remote was set
	url, err := repo.RemoteGet("origin")
	if err != nil {
		t.Errorf("RemoteGet failed: %v", err)
	}
	if url != testURL {
		t.Errorf("expected remote URL %q, got %q", testURL, url)
	}
}

// TestServiceRemoteGet_Success tests RemoteGet retrieves the configured remote.
func TestServiceRemoteGet_Success(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	testURL := "https://github.com/example/repo.git"
	if err := svc.RemoteSet(testURL); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Get remote
	url, err := svc.RemoteGet()
	if err != nil {
		t.Errorf("RemoteGet failed: %v", err)
	}
	if url != testURL {
		t.Errorf("expected remote URL %q, got %q", testURL, url)
	}
}

// TestServiceRemoteGet_NotConfigured tests RemoteGet returns error when no remote.
func TestServiceRemoteGet_NotConfigured(t *testing.T) {
	repo := NewMockRepository()
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// RemoteGet without setting remote should error
	_, err := svc.RemoteGet()
	if err == nil {
		t.Errorf("expected error when remote not configured, got nil")
	}
	if !strings.Contains(err.Error(), "no remote configured") {
		t.Errorf("expected 'no remote configured' in error, got: %v", err)
	}
}
