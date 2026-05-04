package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var importDryRun bool

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Migrate legacy focus sessions to canonical name format",
	Long: `Migrates all legacy-format sessions to the canonical YYYY-MM-DD-HHmm__name format.

Two passes:
  1. Renames legacy git branches in ~/.focus
  2. Renames legacy workspace directories and creates missing branches

Name conversion:
  YYYY-MM-DD--name  →  YYYY-MM-DD-0000__name
  plain-name        →  2000-01-01-0000__plain-name

Use --dry-run to preview changes without modifying anything.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := service.ImportFocuses(importDryRun)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			fmt.Println("Nothing to import — all sessions already in canonical format.")
			return nil
		}
		if importDryRun {
			fmt.Println("Dry run — no changes made:")
		}
		for _, r := range results {
			switch {
			case r.Skipped:
				fmt.Printf("  skip      [%s] %s (%s)\n", r.Source, r.OldName, r.SkipReason)
			case r.SkipReason != "":
				// dir renamed but branch create failed
				fmt.Printf("  partial   [%s] %s → %s (%s)\n", r.Source, r.OldName, r.NewName, r.SkipReason)
			case r.OldName != r.NewName:
				verb := "import"
				if importDryRun {
					verb = "would import"
				}
				fmt.Printf("  %-10s [%s] %s → %s\n", verb, r.Source, r.OldName, r.NewName)
			case r.BranchCreated:
				verb := "new branch"
				if importDryRun {
					verb = "would branch"
				}
				fmt.Printf("  %-10s [%s] %s\n", verb, r.Source, r.OldName)
			}
		}
		return nil
	},
}

func init() {
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "preview changes without modifying anything")
	rootCmd.AddCommand(importCmd)
}
