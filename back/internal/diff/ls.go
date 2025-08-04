package diff

import (
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Path       string
	LastModded int64
}

func listFilesWithModTime(root string) ([]FileInfo, error) {
	var files []FileInfo

	allowedExt := map[string]bool{
		".epub": true,
		".pdf":  true,
		".cbz":  true,
		".cbr":  true,
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if allowedExt[ext] {
				files = append(files, FileInfo{
					Path:       path,
					LastModded: info.ModTime().Unix(),
				})
			}
		}
		return nil
	})

	return files, err
}
