package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote [url]",
	Short: "Get or set the remote URL for backup/sync",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			url, err := service.RemoteGet()
			if err != nil {
				return err
			}
			fmt.Println(url)
			return nil
		}
		if err := service.RemoteSet(args[0]); err != nil {
			return err
		}
		fmt.Printf("Remote set to %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(remoteCmd)
}
