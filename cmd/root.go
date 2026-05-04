package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zadewu/focus/internal/adapters/ui"
	"github.com/zadewu/focus/internal/domain"
)

var service *domain.FocusService

func SetService(svc *domain.FocusService) {
	service = svc
}

var rootCmd = &cobra.Command{
	Use:   "focus",
	Short: "Manage your work sessions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ensureInit()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus()
	},
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func focusRepoDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".focus"), nil
}

func ensureInit() error {
	dir, err := focusRepoDir()
	if err != nil {
		return err
	}

	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create focus dir: %w", err)
		}
		if _, err := runGit(dir, "init"); err != nil {
			return fmt.Errorf("init focus repo: %w", err)
		}
	}

	wsRoot, err := workspaceRoot(dir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(wsRoot, 0o755); err != nil {
		return fmt.Errorf("create workspace root: %w", err)
	}
	return nil
}

// workspaceRoot returns the configured workspace root or the default ~/focus-workspaces.
// Uses git config directly to avoid adapter dependency from cmd.
func workspaceRoot(focusDir string) (string, error) {
	c := exec.Command("git", "config", "--local", "focus.workspace-root")
	c.Dir = focusDir
	out, err := c.Output()
	if err == nil {
		if v := strings.TrimSpace(string(out)); v != "" {
			return v, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, "focus-workspaces"), nil
}

func runGit(dir string, args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func runStatus() error {
	current, wsPath, notes, err := service.Status()
	if err != nil {
		return err
	}
	if current == "" {
		fmt.Println("No active focus. Run: focus new <name>")
		return nil
	}
	ui.PrintStatus(current, wsPath, notes)
	return nil
}
