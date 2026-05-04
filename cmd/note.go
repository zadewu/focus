package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [message]",
	Short: "Add a note to the current focus",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var msg string
		var err error
		if len(args) == 1 {
			msg = strings.TrimSpace(args[0])
		} else {
			msg, err = openEditor("")
			if err != nil {
				return fmt.Errorf("failed to open editor: %w\nSet $EDITOR or use: focus note \"<message>\"", err)
			}
			msg = strings.TrimSpace(msg)
		}
		if msg == "" {
			return fmt.Errorf("note is empty, aborting")
		}
		if err := service.AddNote(msg); err != nil {
			return err
		}
		fmt.Println("Note added")
		return nil
	},
}

// openEditor opens the user's $EDITOR (fallback: vi) with an optional initial string,
// and returns the file contents after the editor exits.
func openEditor(initial string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	f, err := os.CreateTemp("", "focus-note-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	if initial != "" {
		_, _ = f.WriteString(initial)
	}
	f.Close()

	c := exec.Command(editor, f.Name())
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return "", err
	}
	data, err := os.ReadFile(f.Name())
	return string(data), err
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
