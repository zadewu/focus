package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new focus session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fullName, wsPath, err := service.NewFocus(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Created:   %s\nWorkspace: %s\n", fullName, wsPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
