package config

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/zadewu/focus/internal/domain"
)

var _ domain.ConfigStore = (*GitConfigStore)(nil)

type GitConfigStore struct {
	dir string
}

func NewGitConfigStore(dir string) *GitConfigStore {
	return &GitConfigStore{dir: dir}
}

func (s *GitConfigStore) Get(key string) (string, error) {
	c := exec.Command("git", "config", "--local", "focus."+key)
	c.Dir = s.dir
	out, err := c.Output()
	if err != nil {
		// Exit code 1 means key not set — treat as empty, not an error.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("config get %q: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *GitConfigStore) Set(key, value string) error {
	c := exec.Command("git", "config", "--local", "focus."+key, value)
	c.Dir = s.dir
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("config set %q: %s: %w", key, strings.TrimSpace(string(out)), err)
	}
	return nil
}
