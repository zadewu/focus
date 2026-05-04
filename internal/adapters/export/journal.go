package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zadewu/focus/internal/domain"
)

// updateJournal appends or replaces a focus section in today's Obsidian daily note.
// It silently skips if the journal file does not exist or there are no notes today.
func updateJournal(vaultPath, pattern, focusName, relFocusPath string, notes []domain.Note) error {
	if pattern == "" {
		pattern = "01 Daily/{YYYY}/{MM}/{YYYY}-{MM}-{DD}"
	}
	today := time.Now()
	journalPath := filepath.Join(vaultPath, expandJournalPattern(pattern, today))
	if _, err := os.Stat(journalPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: journal not found at %s — skipping\n", journalPath)
		return nil
	}
	var todayNotes []domain.Note
	todayStr := today.Format("2006-01-02")
	for _, n := range notes {
		if n.Timestamp.Format("2006-01-02") == todayStr {
			todayNotes = append(todayNotes, n)
		}
	}
	if len(todayNotes) == 0 {
		return nil
	}
	data, err := os.ReadFile(journalPath)
	if err != nil {
		return fmt.Errorf("read journal: %w", err)
	}
	entry := buildJournalEntry(focusName, relFocusPath, todayNotes)
	content := upsertFocusSection(string(data), focusName, relFocusPath, entry)
	return os.WriteFile(journalPath, []byte(content), 0o644)
}

func expandJournalPattern(pattern string, t time.Time) string {
	r := strings.NewReplacer(
		"{YYYY}", t.Format("2006"),
		"{MM}", t.Format("01"),
		"{DD}", t.Format("02"),
	)
	return r.Replace(pattern) + ".md"
}

func buildJournalEntry(name, relPath string, notes []domain.Note) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "- [[%s|%s]] — %d notes today\n", relPath, name, len(notes))
	for _, n := range notes {
		fmt.Fprintf(&sb, "  - %s — %s\n", n.Timestamp.Format("15:04"), n.Message)
	}
	return sb.String()
}

// upsertFocusSection inserts or replaces the focus link block under the "## Focus" heading.
func upsertFocusSection(content, focusName, relFocusPath, entry string) string {
	linkPrefix := fmt.Sprintf("- [[%s|%s]]", relFocusPath, focusName)
	const header = "## Focus"
	idx := strings.Index(content, header)
	if idx == -1 {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		return content + "\n" + header + "\n\n" + entry
	}
	after := content[idx+len(header):]
	nextIdx := strings.Index(after, "\n## ")
	var section, tail string
	if nextIdx == -1 {
		section, tail = after, ""
	} else {
		section, tail = after[:nextIdx+1], after[nextIdx+1:]
	}
	if strings.Contains(section, linkPrefix) {
		lines := strings.Split(section, "\n")
		var out []string
		skip := false
		for i, line := range lines {
			if strings.HasPrefix(line, linkPrefix) {
				out = append(out, strings.Split(strings.TrimRight(entry, "\n"), "\n")...)
				skip = true
				continue
			}
			if skip && i < len(lines)-1 && (strings.HasPrefix(line, "  ") || line == "") {
				continue
			}
			skip = false
			out = append(out, line)
		}
		section = strings.Join(out, "\n")
	} else {
		section = strings.TrimRight(section, "\n") + "\n" + entry
	}
	return content[:idx+len(header)] + section + tail
}
