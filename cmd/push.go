package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push all sessions to the configured remote",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, err := service.Push()
		if err != nil {
			return err
		}
		fmt.Printf("Pushed to %s\n", url)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
