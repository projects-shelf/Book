package diff

import (
	"back/database"
	"database/sql"
)

type DiffResult struct {
	Added   []string
	Updated []string
	Deleted []string
}

func Diff(db *sql.DB) (DiffResult, error) {
	dbFiles, err := database.GetPathAndLastModdedList(db)
	if err != nil {
		return DiffResult{}, err
	}

	fsFiles, err := listFilesWithModTime("/books/")
	if err != nil {
		return DiffResult{}, err
	}

	dbMap := make(map[string]int64)
	for _, f := range dbFiles {
		dbMap[f.Path] = f.LastModded
	}

	fsMap := make(map[string]int64)
	for _, f := range fsFiles {
		fsMap[f.Path] = f.LastModded
	}

	var added, updated, deleted []string

	for path := range fsMap {
		if _, ok := dbMap[path]; !ok {
			added = append(added, path)
		}
	}

	for path := range dbMap {
		if _, ok := fsMap[path]; !ok {
			deleted = append(deleted, path)
		}
	}

	for path, dbMod := range dbMap {
		if fsMod, ok := fsMap[path]; ok && fsMod != dbMod {
			updated = append(updated, path)
		}
	}

	return DiffResult{
		Added:   added,
		Updated: updated,
		Deleted: deleted,
	}, nil
}
