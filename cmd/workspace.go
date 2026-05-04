package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace [name]",
	Short: "Print workspace path for a focus",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) == 1 {
			name = args[0]
		}
		path, err := service.WorkspacePath(name)
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
}
