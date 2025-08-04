package scan

import (
	"back/database"
	"database/sql"
	"os"
)

func scanDelete(path string, bookDB, keywordDB *sql.DB) error {
	book, err := database.GetBookByPath(bookDB, path)

	if err != nil {
		return err
	}

	deleteCoverFile(book.CoverPath)

	err = database.DeleteBookByPath(bookDB, path)
	if err != nil {
		return err
	}

	err = database.DeleteKeywordsByPath(keywordDB, path)
	if err != nil {
		return err
	}

	return nil
}

func deleteCoverFile(coverPath string) error {
	if coverPath == "" {
		return nil
	}
	return os.Remove(coverPath)
}
