package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

func OpenKeywordDB() *sql.DB {
	db, err := sql.Open("sqlite", "/db/keyword.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS book_keywords (
			path TEXT NOT NULL,
			keyword TEXT NOT NULL,
			PRIMARY KEY (path, keyword)
		);
	`)
	if err != nil {
		panic(err)
	}

	return db
}

func AddKeyword(db *sql.DB, path string, keyword string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO book_keywords (path, keyword)
		VALUES (?, ?)
	`, path, keyword)
	if err != nil {
		return fmt.Errorf("failed to add keyword: %w", err)
	}
	return nil
}

func DeleteKeywordsByPath(db *sql.DB, path string) error {
	_, err := db.Exec(`DELETE FROM book_keywords WHERE path = ?`, path)
	return err
}

func FindPathsByKeywords(db *sql.DB, keywordList []string) ([]string, error) {
	if len(keywordList) == 0 {
		return []string{}, nil
	}

	placeholders := strings.Repeat("?,", len(keywordList))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]interface{}, len(keywordList)+1)
	for i, kw := range keywordList {
		args[i] = kw
	}
	args[len(keywordList)] = len(keywordList)

	query := fmt.Sprintf(`
		SELECT path
		FROM book_keywords
		WHERE keyword IN (%s)
		GROUP BY path
		HAVING COUNT(DISTINCT keyword) = ?;
	`, placeholders)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	return paths, nil
}
