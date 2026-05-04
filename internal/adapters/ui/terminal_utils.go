package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
)

// getTerminalWidth returns the terminal column count, defaulting to 80.
func getTerminalWidth() int {
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

// wordWrap breaks s into lines of at most maxWidth visible characters.
// Existing newlines in s are preserved; continuation lines are prefixed
// with indent spaces for alignment.
func wordWrap(s string, maxWidth, indent int) string {
	if maxWidth <= 0 {
		return s
	}
	pad := strings.Repeat(" ", indent)
	segments := strings.Split(s, "\n")
	wrapped := make([]string, 0, len(segments))
	for _, seg := range segments {
		wrapped = append(wrapped, wrapSegment(seg, pad, maxWidth))
	}
	return strings.Join(wrapped, "\n"+pad)
}

// wrapSegment wraps a single line (no embedded newlines) at maxWidth characters.
func wrapSegment(s, pad string, maxWidth int) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}
	var lines []string
	cur := ""
	for _, w := range words {
		switch {
		case cur == "":
			cur = w
		case len(cur)+1+len(w) <= maxWidth:
			cur += " " + w
		default:
			lines = append(lines, cur)
			cur = w
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return strings.Join(lines, "\n"+pad)
}
