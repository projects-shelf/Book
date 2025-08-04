package scan

import (
	"back/database"
	"back/internal/cover"
	"back/internal/meta"
	"database/sql"
	"fmt"
)

func scanUpdate(path string, bookDB, keywordDB *sql.DB) error {
	book, err := database.GetBookByPath(bookDB, path)

	if err != nil {
		return err
	}

	err = cover.ExtractCover(path, book.CoverPath, book.Type)
	if err != nil {
		fmt.Println(err)
	}

	title, keywords, last_modded, err := meta.ExtractMeta(path, book.Type)
	if err != nil {
		return err
	}

	database.UpdateBookTitleAndModTime(bookDB, path, title, last_modded)

	database.DeleteKeywordsByPath(keywordDB, path)
	for _, keyword := range keywords {
		err = database.AddKeyword(keywordDB, path, keyword)
		if err != nil {
			return err
		}
	}

	return nil
}
