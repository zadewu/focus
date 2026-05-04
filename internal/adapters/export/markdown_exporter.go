package export

import (
	"fmt"
	"os"
	"strings"

	"github.com/zadewu/focus/internal/domain"
)

var _ domain.Exporter = (*MarkdownExporter)(nil)

type MarkdownExporter struct{ outDir string }

func NewMarkdown(outDir string) *MarkdownExporter {
	return &MarkdownExporter{outDir: outDir}
}

func (e *MarkdownExporter) Export(focus domain.Focus, notes []domain.Note, files []domain.File) error {
	content := buildMarkdown(focus, notes, files)
	outPath := e.outDir + "/" + focus.Name + ".md"
	if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write export: %w", err)
	}
	fmt.Printf("Exported: %s\n", outPath)
	return nil
}

// buildMarkdown assembles a markdown document from focus metadata, notes, and workspace files.
func buildMarkdown(focus domain.Focus, notes []domain.Note, files []domain.File) string {
	var sb strings.Builder
	sb.WriteString("# Focus: " + focus.Name + "\n\n")
	status := "Active"
	if focus.Archived {
		status = "Archived"
	}
	if !focus.CreatedAt.IsZero() {
		sb.WriteString("**Created:** " + focus.CreatedAt.Format("2006-01-02 15:04") + "\n")
	}
	if len(notes) > 0 && !notes[0].Timestamp.IsZero() {
		sb.WriteString("**Last Updated:** " + notes[0].Timestamp.Format("2006-01-02 15:04") + "\n")
	}
	sb.WriteString("**Status:** " + status + "\n\n")
	if len(notes) > 0 {
		sb.WriteString("## Notes\n\n")
		for _, n := range notes {
			ts := "unknown"
			if !n.Timestamp.IsZero() {
				ts = n.Timestamp.Format("2006-01-02 15:04")
			}
			sb.WriteString("### " + ts + "\n")
			sb.WriteString(n.Message + "\n\n")
		}
		sb.WriteString("---\n\n")
	}
	var mdFiles []domain.File
	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name), ".md") {
			mdFiles = append(mdFiles, f)
		}
	}
	if len(mdFiles) > 0 {
		sb.WriteString("## Workspace Files\n\n")
		for _, f := range mdFiles {
			sb.WriteString("### " + f.Name + "\n\n")
			sb.WriteString(string(f.Content) + "\n\n")
		}
	}
	return sb.String()
}
