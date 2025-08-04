package meta

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CBRMeta holds extracted metadata
type CBRMeta struct {
	Title    string
	Keywords []string
	ModTime  time.Time
}

// TODO
func extractCBRMeta(path string) (string, []string, int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", nil, 0, fmt.Errorf("failed to stat file: %w", err)
	}
	modTime := info.ModTime().Unix()

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	var keywords []string

	return title, keywords, modTime, nil
}
