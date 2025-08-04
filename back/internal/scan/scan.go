package scan

import (
	"back/internal/diff"
	"database/sql"
	"fmt"
)

func Scan(bookDB, keywordDB *sql.DB) error {
	diffResult, err := diff.Diff(bookDB)
	if err != nil {
		return err
	}

	for _, path := range diffResult.Added {
		err := scanAdd(path, bookDB, keywordDB)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, path := range diffResult.Updated {
		err := scanUpdate(path, bookDB, keywordDB)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, path := range diffResult.Deleted {
		err := scanDelete(path, bookDB, keywordDB)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
