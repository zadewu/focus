package git

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// setupLocalRepo creates a temporary local git repository with initial commit.
// Returns the repo path and cleanup function.
func setupLocalRepo(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "focus-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init repo: %v", err)
	}

	// Set git user for commits
	gitCmd := func(args ...string) error {
		c := exec.Command("git", args...)
		c.Dir = tmpDir
		return c.Run()
	}

	if err := gitCmd("config", "user.email", "test@example.com"); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to set git user email: %v", err)
	}
	if err := gitCmd("config", "user.name", "Test User"); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to set git user name: %v", err)
	}

	// Create initial commit
	if err := gitCmd("commit", "--allow-empty", "-m", "init"); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to create initial commit: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return tmpDir, cleanup
}

// setupBareRepo creates a temporary bare git repository.
// Returns the repo path and cleanup function.
func setupBareRepo(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "focus-bare-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cmd := exec.Command("git", "init", "--bare", tmpDir)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init bare repo: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return tmpDir, cleanup
}

// gitRun executes a git command in a given directory and returns the output.
func gitRun(t *testing.T, dir string, args ...string) string {
	c := exec.Command("git", args...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	if err != nil {
		t.Logf("git command failed: git %s (in %s): %v\noutput: %s", strings.Join(args, " "), dir, err, string(out))
	}
	return strings.TrimSpace(string(out))
}

// TestRemoteSet_Add tests setting a remote on a fresh repo adds it.
func TestRemoteSet_Add(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir, cleanupBare := setupBareRepo(t)
	defer cleanupBare()

	repo := NewRepository(localDir)

	// RemoteSet should add a new remote
	err := repo.RemoteSet("origin", bareDir)
	if err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Verify remote was added
	out := gitRun(t, localDir, "remote", "-v")
	if !strings.Contains(out, "origin") {
		t.Errorf("expected 'origin' in remote list, got: %s", out)
	}
	if !strings.Contains(out, bareDir) {
		t.Errorf("expected bare repo path in remote list, got: %s", out)
	}
}

// TestRemoteSet_Update tests updating an existing remote URL.
func TestRemoteSet_Update(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir1, cleanupBare1 := setupBareRepo(t)
	defer cleanupBare1()
	bareDir2, cleanupBare2 := setupBareRepo(t)
	defer cleanupBare2()

	repo := NewRepository(localDir)

	// Set initial remote
	if err := repo.RemoteSet("origin", bareDir1); err != nil {
		t.Fatalf("RemoteSet add failed: %v", err)
	}

	// Update remote URL
	if err := repo.RemoteSet("origin", bareDir2); err != nil {
		t.Fatalf("RemoteSet update failed: %v", err)
	}

	// Verify remote was updated
	out := gitRun(t, localDir, "remote", "get-url", "origin")
	if out != bareDir2 {
		t.Errorf("expected remote URL %q, got %q", bareDir2, out)
	}
}

// TestRemoteGet_Missing tests getting a non-existent remote returns error.
func TestRemoteGet_Missing(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	repo := NewRepository(localDir)

	_, err := repo.RemoteGet("nonexistent")
	if err == nil {
		t.Errorf("expected error for missing remote, got nil")
	}
}

// TestRemoteGet_Present tests getting an existing remote returns correct URL.
func TestRemoteGet_Present(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir, cleanupBare := setupBareRepo(t)
	defer cleanupBare()

	repo := NewRepository(localDir)

	// Set remote
	if err := repo.RemoteSet("origin", bareDir); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Get remote
	url, err := repo.RemoteGet("origin")
	if err != nil {
		t.Fatalf("RemoteGet failed: %v", err)
	}

	if url != bareDir {
		t.Errorf("expected URL %q, got %q", bareDir, url)
	}
}

// TestPushAll tests pushing all branches and tags to remote.
func TestPushAll(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir, cleanupBare := setupBareRepo(t)
	defer cleanupBare()

	repo := NewRepository(localDir)

	// Set up remote
	if err := repo.RemoteSet("origin", bareDir); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Create additional branches and commits
	gitRun(t, localDir, "checkout", "-b", "feature1")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "feature commit")

	gitRun(t, localDir, "checkout", "-b", "feature2")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "another feature")

	// Create a tag
	gitRun(t, localDir, "tag", "v1.0")

	// Push all
	if err := repo.PushAll("origin"); err != nil {
		t.Fatalf("PushAll failed: %v", err)
	}

	// Verify branches exist in bare repo
	branches := gitRun(t, bareDir, "branch", "-a")
	if !strings.Contains(branches, "master") && !strings.Contains(branches, "main") {
		t.Errorf("expected default branch in bare repo, got: %s", branches)
	}
	if !strings.Contains(branches, "feature1") {
		t.Errorf("expected feature1 branch in bare repo, got: %s", branches)
	}
	if !strings.Contains(branches, "feature2") {
		t.Errorf("expected feature2 branch in bare repo, got: %s", branches)
	}

	// Verify tag exists in bare repo
	tags := gitRun(t, bareDir, "tag")
	if !strings.Contains(tags, "v1.0") {
		t.Errorf("expected v1.0 tag in bare repo, got: %s", tags)
	}
}

// TestPushAll_NoRemote tests PushAll returns error when remote doesn't exist.
func TestPushAll_NoRemote(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	repo := NewRepository(localDir)

	// Try to push to non-existent remote
	err := repo.PushAll("origin")
	if err == nil {
		t.Errorf("expected error when pushing to non-existent remote")
	}
}

// TestFetchAll tests fetching from remote.
func TestFetchAll(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir, cleanupBare := setupBareRepo(t)
	defer cleanupBare()

	repo := NewRepository(localDir)

	// Set up remote and push initial state
	if err := repo.RemoteSet("origin", bareDir); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}
	if err := repo.PushAll("origin"); err != nil {
		t.Fatalf("PushAll failed: %v", err)
	}

	// Create new branches in local repo
	gitRun(t, localDir, "checkout", "-b", "new-feature")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "new feature")
	gitRun(t, localDir, "push", "-u", "origin", "new-feature")

	// Create second local clone from bare
	clone2Dir, cleanup2 := setupLocalRepo(t)
	defer cleanup2()

	gitRun(t, clone2Dir, "remote", "add", "origin", bareDir)
	gitRun(t, clone2Dir, "config", "user.email", "test@example.com")
	gitRun(t, clone2Dir, "config", "user.name", "Test User")

	repo2 := NewRepository(clone2Dir)

	// Fetch from remote should get new-feature branch
	if err := repo2.FetchAll("origin"); err != nil {
		t.Fatalf("FetchAll failed: %v", err)
	}

	// Verify remote branch exists
	remoteBranches := gitRun(t, clone2Dir, "branch", "-r")
	if !strings.Contains(remoteBranches, "origin/new-feature") {
		t.Errorf("expected origin/new-feature after FetchAll, got: %s", remoteBranches)
	}
}

// TestFetchAll_NoRemote tests FetchAll returns an error when the named remote doesn't exist.
func TestFetchAll_NoRemote(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	repo := NewRepository(localDir)

	err := repo.FetchAll("origin")
	if err == nil {
		t.Errorf("FetchAll with no remote configured should return error, got nil")
	}
}

// TestCheckoutRemoteBranches tests creating local branches from remote refs.
func TestCheckoutRemoteBranches(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir, cleanupBare := setupBareRepo(t)
	defer cleanupBare()

	repo := NewRepository(localDir)

	// Set up remote and initial push
	if err := repo.RemoteSet("origin", bareDir); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Create and push multiple branches
	gitRun(t, localDir, "checkout", "-b", "develop")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "develop init")
	gitRun(t, localDir, "checkout", "-b", "feature/auth")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "auth feature")
	gitRun(t, localDir, "checkout", "-b", "feature/api")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "api feature")

	if err := repo.PushAll("origin"); err != nil {
		t.Fatalf("PushAll failed: %v", err)
	}

	// Create second local repo cloned from bare
	clone2Dir, cleanup2 := setupLocalRepo(t)
	defer cleanup2()

	gitRun(t, clone2Dir, "remote", "add", "origin", bareDir)
	gitRun(t, clone2Dir, "config", "user.email", "test@example.com")
	gitRun(t, clone2Dir, "config", "user.name", "Test User")

	repo2 := NewRepository(clone2Dir)

	// Fetch and then checkout remote branches
	if err := repo2.FetchAll("origin"); err != nil {
		t.Fatalf("FetchAll failed: %v", err)
	}

	if err := repo2.CheckoutRemoteBranches("origin"); err != nil {
		t.Fatalf("CheckoutRemoteBranches failed: %v", err)
	}

	// Verify all remote branches are now local
	branches := gitRun(t, clone2Dir, "branch")

	// Find which default branch was created
	hasDefault := strings.Contains(branches, "master") || strings.Contains(branches, "main")
	if !hasDefault {
		t.Errorf("expected master or main branch, got: %s", branches)
	}

	for _, expected := range []string{"develop", "feature/auth", "feature/api"} {
		if !strings.Contains(branches, expected) {
			t.Errorf("expected branch %q in: %s", expected, branches)
		}
	}
}

// TestCheckoutRemoteBranches_SkipsExisting tests CheckoutRemoteBranches doesn't re-checkout existing branches.
func TestCheckoutRemoteBranches_SkipsExisting(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	bareDir, cleanupBare := setupBareRepo(t)
	defer cleanupBare()

	repo := NewRepository(localDir)

	// Set up remote
	if err := repo.RemoteSet("origin", bareDir); err != nil {
		t.Fatalf("RemoteSet failed: %v", err)
	}

	// Create a branch and push
	gitRun(t, localDir, "checkout", "-b", "existing-feature")
	gitRun(t, localDir, "commit", "--allow-empty", "-m", "existing feature")

	if err := repo.PushAll("origin"); err != nil {
		t.Fatalf("PushAll failed: %v", err)
	}

	// Call CheckoutRemoteBranches on the same repo
	// It should skip existing-feature since it's already local
	if err := repo.CheckoutRemoteBranches("origin"); err != nil {
		t.Fatalf("CheckoutRemoteBranches failed: %v", err)
	}

	// Verify existing-feature is still there
	branches := gitRun(t, localDir, "branch")
	if !strings.Contains(branches, "existing-feature") {
		t.Errorf("expected existing-feature to still exist, got: %s", branches)
	}
}

// TestCheckoutRemoteBranches_NoRemote tests CheckoutRemoteBranches succeeds with no remote branches.
// It gracefully handles repos with no remote branches and just returns nil.
func TestCheckoutRemoteBranches_NoRemote(t *testing.T) {
	localDir, cleanup := setupLocalRepo(t)
	defer cleanup()

	repo := NewRepository(localDir)

	// CheckoutRemoteBranches with no remotes should succeed (no-op)
	err := repo.CheckoutRemoteBranches("origin")
	if err != nil {
		t.Errorf("CheckoutRemoteBranches with no remotes should succeed, got error: %v", err)
	}
}
