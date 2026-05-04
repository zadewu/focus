package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/zadewu/focus/cmd"
	configadapter "github.com/zadewu/focus/internal/adapters/config"
	gitadapter "github.com/zadewu/focus/internal/adapters/git"
	wsadapter "github.com/zadewu/focus/internal/adapters/workspace"
	"github.com/zadewu/focus/internal/domain"
)

// version is set at build time via -ldflags "-X main.version=vX.Y.Z"
var version = "dev"

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	focusDir := filepath.Join(home, ".focus")
	repo := gitadapter.NewRepository(focusDir)
	cfg := configadapter.NewGitConfigStore(focusDir)

	// Resolve workspace root from config, falling back to ~/focus-workspaces.
	wsRoot, err := cfg.Get("workspace-root")
	if err != nil || wsRoot == "" {
		wsRoot = filepath.Join(home, "focus-workspaces")
	}
	ws := wsadapter.NewFSWorkspaceStore(wsRoot)

	svc := domain.NewFocusService(repo, cfg, ws)
	cmd.SetService(svc)
	cmd.Execute(version)
}
