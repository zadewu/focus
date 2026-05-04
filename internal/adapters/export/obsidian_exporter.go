package export

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zadewu/focus/internal/domain"
)

var _ domain.Exporter = (*ObsidianExporter)(nil)

type ObsidianExporter struct {
	vaultPath      string
	workspaceRoot  string
	journalPattern string
}

func NewObsidian(vaultPath, workspaceRoot, journalPattern string) *ObsidianExporter {
	return &ObsidianExporter{
		vaultPath:      expandHome(vaultPath),
		workspaceRoot:  expandHome(workspaceRoot),
		journalPattern: journalPattern,
	}
}

func (e *ObsidianExporter) Export(focus domain.Focus, notes []domain.Note, files []domain.File) error {
	focusFilename := focus.Name + ".md"
	focusVaultDir := filepath.Join(e.vaultPath, "Focus")
	if err := os.MkdirAll(focusVaultDir, 0o755); err != nil {
		return fmt.Errorf("create vault Focus dir: %w", err)
	}
	content := buildMarkdown(focus, notes, files)
	focusFilePath := filepath.Join(focusVaultDir, focusFilename)
	if err := os.WriteFile(focusFilePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write vault file: %w", err)
	}
	fmt.Printf("Written: %s\n", focusFilePath)

	wsDir := filepath.Join(e.workspaceRoot, focus.Name)
	if err := zipWorkspaceFiles(wsDir, focus.Name, e.vaultPath); err != nil {
		fmt.Fprintf(os.Stderr, "warning: zip: %v\n", err)
	}

	relFocusPath := "Focus/" + strings.TrimSuffix(focusFilename, ".md")
	if err := updateJournal(e.vaultPath, e.journalPattern, focus.Name, relFocusPath, notes); err != nil {
		fmt.Fprintf(os.Stderr, "warning: journal update: %v\n", err)
	}
	return nil
}

func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}

// zipWorkspaceFiles archives non-markdown files from wsDir into vault/Focus/attachments/<name>.zip.
func zipWorkspaceFiles(wsDir, name, vaultPath string) error {
	if _, err := os.Stat(wsDir); os.IsNotExist(err) {
		return nil
	}
	var filesToZip []string
	_ = filepath.Walk(wsDir, func(p string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(p), ".md") {
			filesToZip = append(filesToZip, p)
		}
		return nil
	})
	if len(filesToZip) == 0 {
		return nil
	}
	attDir := filepath.Join(vaultPath, "Focus", "attachments")
	if err := os.MkdirAll(attDir, 0o755); err != nil {
		return err
	}
	zipPath := filepath.Join(attDir, name+".zip")
	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()
	w := zip.NewWriter(zf)
	defer w.Close()
	for _, p := range filesToZip {
		rel, _ := filepath.Rel(wsDir, p)
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		zh, err := zip.FileInfoHeader(info)
		if err != nil {
			continue
		}
		zh.Name = filepath.ToSlash(rel)
		zh.Method = zip.Deflate
		fw, err := w.CreateHeader(zh)
		if err != nil {
			continue
		}
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		_, _ = io.Copy(fw, f)
		f.Close()
	}
	fmt.Printf("Zipped:  %s\n", zipPath)
	return nil
}
