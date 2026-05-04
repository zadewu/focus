package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to an existing focus",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, err := service.SwitchFocus(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Switched to: %s\nWorkspace:   %s\n", args[0], wsPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
