package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pullRestore bool

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fetch sessions from the configured remote",
	Long: `Fetches all branches from the remote.

Use --restore on a fresh machine to create local tracking branches
for every remote branch (useful after migrating to a new machine).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, err := service.RemoteGet()
		if err != nil {
			return err
		}
		fmt.Printf("Fetching from %s...\n", url)
		if err := service.Pull(pullRestore); err != nil {
			return err
		}
		if pullRestore {
			fmt.Println("Done. Local branches created from remote.")
		} else {
			fmt.Println("Done. Run 'focus pull --restore' to create local branches.")
		}
		return nil
	},
}

func init() {
	pullCmd.Flags().BoolVar(&pullRestore, "restore", false, "create local tracking branches for all remote branches")
	rootCmd.AddCommand(pullCmd)
}
