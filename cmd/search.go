package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	ui "github.com/zadewu/focus/internal/adapters/ui"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search notes across all active and archived sessions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := args[0]
		results, err := service.SearchNotes(keyword)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			fmt.Printf("No notes matching %q.\n", keyword)
			return nil
		}
		ui.PrintSearchResults(keyword, results)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
