package meta

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// PDFMeta holds extracted metadata
type PDFMeta struct {
	Title    string
	Keywords []string
	ModTime  time.Time
}

// ExtractPDFMeta extracts metadata and last modified time from the given PDF file
func extractPDFMeta(path string) (string, []string, int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", nil, 0, fmt.Errorf("failed to stat file: %w", err)
	}
	modTime := info.ModTime().Unix()

	cmd := exec.Command("pdfinfo", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", nil, 0, fmt.Errorf("failed to run pdfinfo: %w", err)
	}

	lines := strings.Split(out.String(), "\n")
	meta := make(map[string]string)
	for _, line := range lines {
		if sep := strings.Index(line, ":"); sep != -1 {
			key := strings.TrimSpace(line[:sep])
			value := strings.TrimSpace(line[sep+1:])
			meta[key] = value
		}
	}

	title := meta["Title"]
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	var keywords []string
	if author := meta["Author"]; author != "" {
		keywords = append(keywords, author)
	}
	if kw := meta["Keywords"]; kw != "" {
		keywords = append(keywords, strings.Split(kw, ",")...)
	}

	return title, keywords, modTime, nil
}
