package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/zadewu/focus/internal/domain"
)

// mutedColor is readable on both dark/transparent and light terminals.
var mutedColor = lipgloss.AdaptiveColor{Light: "#606060", Dark: "#aaaaaa"}

var (
	// ActiveStyle renders active focus names in bold green.
	ActiveStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	// ArchivedStyle renders archived focus names in a muted tone readable on any background.
	ArchivedStyle = lipgloss.NewStyle().Foreground(mutedColor)
	// CurrentMark renders the active-focus indicator arrow.
	CurrentMark = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	// HeaderStyle renders section headers.
	HeaderStyle = lipgloss.NewStyle().Bold(true).Underline(true)
	// DimStyle renders secondary/supplementary text.
	DimStyle = lipgloss.NewStyle().Foreground(mutedColor)
	// NoteStyle renders note message text.
	NoteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	// TimestampStyle renders note timestamps in a muted tone readable on any background.
	TimestampStyle = lipgloss.NewStyle().Foreground(mutedColor)
)

// PrintFocusList renders the full focus list, separating active from archived.
func PrintFocusList(focuses []domain.Focus, current string) {
	var active, archived []domain.Focus
	for _, f := range focuses {
		if f.Archived {
			archived = append(archived, f)
		} else {
			active = append(active, f)
		}
	}
	for _, f := range active {
		if f.Name == current {
			fmt.Println(CurrentMark.Render("▶ ") + ActiveStyle.Render(f.Name))
		} else {
			fmt.Println("  " + f.Name)
		}
	}
	if len(archived) > 0 {
		fmt.Println(DimStyle.Render("── archived ──"))
		for _, f := range archived {
			fmt.Println(ArchivedStyle.Render("  " + f.Name))
		}
	}
}

// logPrefixLen is the visible width of "  HH:MM  " (2 spaces + 5 time chars + 2 spaces).
const logPrefixLen = 9

// PrintLog renders the full note history for a focus.
func PrintLog(focusName string, notes []domain.Note) {
	fmt.Println(HeaderStyle.Render(focusName))
	fmt.Println()
	termW := getTerminalWidth()
	for _, n := range notes {
		ts := FormatTimestamp(n.Timestamp)
		msg := wordWrap(n.Message, termW-logPrefixLen, logPrefixLen)
		fmt.Printf("  %s  %s\n", TimestampStyle.Render(ts), NoteStyle.Render(msg))
	}
}

// PrintStatus renders the current focus name, workspace path, and up to 5 recent notes.
func PrintStatus(current, wsPath string, notes []domain.Note) {
	fmt.Println(CurrentMark.Render("▶ ") + ActiveStyle.Render(current))
	fmt.Println(DimStyle.Render("   Workspace: " + wsPath))
	fmt.Println()
	fmt.Println(HeaderStyle.Render("Recent notes:"))
	if len(notes) == 0 {
		fmt.Println(DimStyle.Render("  (none yet — run: focus note <message>)"))
		return
	}
	termW := getTerminalWidth()
	limit := min(5, len(notes))
	for _, n := range notes[:limit] {
		ts := FormatTimestamp(n.Timestamp)
		msg := wordWrap(n.Message, termW-logPrefixLen, logPrefixLen)
		fmt.Printf("  %s  %s\n", TimestampStyle.Render(ts), NoteStyle.Render(msg))
	}
}

// PrintSearchResults renders search matches with their focus context.
func PrintSearchResults(keyword string, results []domain.SearchResult) {
	fmt.Printf("Search results for %s:\n\n", HeaderStyle.Render(keyword))
	for _, r := range results {
		focusLabel := ActiveStyle.Render(domain.ExtractShortName(r.Focus.Name))
		if r.Focus.Archived {
			focusLabel = ArchivedStyle.Render(domain.ExtractShortName(r.Focus.Name)) +
				DimStyle.Render(" (archived)")
		}
		ts := FormatDate(r.Note.Timestamp)
		fmt.Printf("  %s  %s  %s\n",
			focusLabel,
			TimestampStyle.Render(ts),
			NoteStyle.Render(r.Note.Message),
		)
	}
}

// FormatTimestamp returns "HH:MM" for inline note display.
func FormatTimestamp(t time.Time) string {
	if t.IsZero() {
		return "??:??"
	}
	return t.Format("15:04")
}

// FormatDate returns "YYYY-MM-DD HH:MM" for log display.
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	return t.Format("2006-01-02 15:04")
}

