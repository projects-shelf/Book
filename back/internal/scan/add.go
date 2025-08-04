package scan

import (
	"back/database"
	"back/internal/cover"
	"back/internal/meta"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func scanAdd(path string, bookDB, keywordDB *sql.DB) error {
	trimmed := strings.TrimPrefix(path, "/book")
	ext := filepath.Ext(trimmed)
	base := strings.TrimSuffix(trimmed, ext)
	coverPath := "/cache/cover" + base + ".webp"
	bookType := detectBookType(path)

	err := cover.ExtractCover(path, coverPath, bookType)
	if err != nil {
		fmt.Println(err)
		coverPath = ""
	}

	title, keywords, last_modded, err := meta.ExtractMeta(path, bookType)
	if err != nil {
		return err
	}

	book := database.BookData{
		Path:            path,
		CoverPath:       coverPath,
		Type:            bookType,
		Title:           title,
		AddedTime:       time.Now().Unix(),
		LastModded:      last_modded,
		LastOpened:      0,
		CurrentPosition: "",
		Progress:        0.0,
	}
	database.AddBook(bookDB, book)

	for _, keyword := range keywords {
		err = database.AddKeyword(keywordDB, path, keyword)
		if err != nil {
			return err
		}
	}

	return nil
}

func detectBookType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".epub":
		return "EPUB"
	case ".pdf":
		return "PDF"
	case ".cbz":
		return "CBZ"
	case ".cbr":
		return "CBR"
	default:
		return "UNKNOWN"
	}
}
