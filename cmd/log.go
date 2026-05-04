package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	ui "github.com/zadewu/focus/internal/adapters/ui"
)

var logCmd = &cobra.Command{
	Use:   "log [name]",
	Short: "Show note history for a focus",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) == 1 {
			name = args[0]
		}
		focusName, notes, err := service.GetLog(name)
		if err != nil {
			return err
		}
		if len(notes) == 0 {
			fmt.Printf("No notes for %s yet.\n", focusName)
			return nil
		}
		ui.PrintLog(focusName, notes)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
