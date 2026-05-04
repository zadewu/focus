package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// validConfigKeys lists the known configuration keys; unknown keys produce a warning but are not blocked.
var validConfigKeys = map[string]bool{
	"workspace-root":           true,
	"obsidian-vault":           true,
	"obsidian-journal-pattern": true,
}

var configCmd = &cobra.Command{
	Use:   "config <key> [value]",
	Short: "Get or set focus configuration",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		if !validConfigKeys[key] {
			fmt.Printf("Warning: unknown config key %q\n", key)
		}
		if len(args) == 1 {
			v, err := service.GetConfig(key)
			if err != nil {
				return err
			}
			fmt.Println(v)
			return nil
		}
		return service.SetConfig(key, args[1])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
