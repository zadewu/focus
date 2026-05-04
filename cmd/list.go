package cmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
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
		// Interactive TUI when connected to a real terminal; plain list when piped.
		if isatty.IsTerminal(os.Stdout.Fd()) {
			selected, err := ui.RunInteractiveList(focuses, current)
			if err != nil {
				return err
			}
			if selected != "" && selected != current {
				wsPath, switchErr := service.SwitchFocus(selected)
				if switchErr != nil {
					return switchErr
				}
				fmt.Printf("Switched to: %s\nWorkspace:   %s\n", selected, wsPath)
			}
			return nil
		}
		ui.PrintFocusList(focuses, current)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
