package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive [name]",
	Short: "Archive a focus session (defaults to current)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		if err := service.ArchiveFocus(name); err != nil {
			return err
		}
		if name == "" {
			fmt.Println("Archived current focus")
		} else {
			fmt.Printf("Archived: %s\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
