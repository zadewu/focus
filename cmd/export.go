package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zadewu/focus/internal/adapters/export"
)

var obsidianFlag bool

var exportCmd = &cobra.Command{
	Use:   "export [name]",
	Short: "Export a focus to markdown (or Obsidian vault)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) == 1 {
			name = args[0]
		}
		if obsidianFlag {
			vaultPath, err := service.GetConfig("obsidian-vault")
			if err != nil || vaultPath == "" {
				return fmt.Errorf("obsidian vault not configured\nRun: focus config obsidian-vault <path>")
			}
			wsRoot, err := service.WorkspaceRoot()
			if err != nil {
				wsRoot = defaultWorkspaceRoot()
			}
			pattern, _ := service.GetConfig("obsidian-journal-pattern")
			if pattern == "" {
				pattern = "01 Daily/{YYYY}/{MM}/{YYYY}-{MM}-{DD}"
			}
			exporter := export.NewObsidian(expandHomeCmd(vaultPath), wsRoot, pattern)
			return service.Export(name, exporter)
		}
		exporter := export.NewMarkdown(".")
		return service.Export(name, exporter)
	},
}

func defaultWorkspaceRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "focus-workspaces"
	}
	return filepath.Join(home, "focus-workspaces")
}

func expandHomeCmd(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}

func init() {
	exportCmd.Flags().BoolVar(&obsidianFlag, "obsidian", false, "Export to Obsidian vault")
	rootCmd.AddCommand(exportCmd)
}
