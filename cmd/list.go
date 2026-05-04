package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	ui "github.com/zadewu/focus/internal/adapters/ui"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all focus sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		focuses, current, err := service.ListFocuses()
		if err != nil {
			return err
		}
		if len(focuses) == 0 {
			fmt.Println("No focuses yet. Run: focus new <name>")
			return nil
		}
		ui.PrintFocusList(focuses, current)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
