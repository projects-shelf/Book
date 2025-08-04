package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type BookData struct {
	Path            string  `json:"path"`
	CoverPath       string  `json:"cover_path"`
	Type            string  `json:"type"`
	Title           string  `json:"title"`
	AddedTime       int64   `json:"added_time"`
	LastModded      int64   `json:"last_modded"`
	LastOpened      int64   `json:"last_opened"`
	CurrentPosition string  `json:"current_position"`
	Progress        float64 `json:"progress"`
}

func OpenBookDB() *sql.DB {
	db, err := sql.Open("sqlite", "/db/book.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			path TEXT PRIMARY KEY,
			cover_path TEXT,
			type TEXT,
			title TEXT,
			added_time INTEGER,
			last_modded INTEGER,
			last_opened INTEGER,
			current_position TEXT,
			progress REAL
		);
	`)
	if err != nil {
		panic(err)
	}

	return db
}

type PathModded struct {
	Path       string
	LastModded int64
}

func GetPathAndLastModdedList(db *sql.DB) ([]PathModded, error) {
	rows, err := db.Query(`SELECT path, last_modded FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PathModded
	for rows.Next() {
		var pm PathModded
		if err := rows.Scan(&pm.Path, &pm.LastModded); err != nil {
			return nil, err
		}
		results = append(results, pm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func AddBook(db *sql.DB, book BookData) error {
	_, err := db.Exec(`
		INSERT INTO books (
			path,
			cover_path,
			type,
			title,
			added_time,
			last_modded,
			last_opened,
			current_position,
			progress
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		book.Path,
		book.CoverPath,
		book.Type,
		book.Title,
		book.AddedTime,
		book.LastModded,
		book.LastOpened,
		book.CurrentPosition,
		book.Progress,
	)
	return err
}

func GetBookByPath(db *sql.DB, path string) (*BookData, error) {
	row := db.QueryRow(`
		SELECT path, cover_path, type, title, added_time,
		       last_modded, last_opened, current_position, progress
		FROM books
		WHERE path = ?`, path)

	var book BookData
	err := row.Scan(
		&book.Path,
		&book.CoverPath,
		&book.Type,
		&book.Title,
		&book.AddedTime,
		&book.LastModded,
		&book.LastOpened,
		&book.CurrentPosition,
		&book.Progress,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid path")
		}
		return nil, err
	}

	return &book, nil
}

func UpdateBookTitleAndModTime(db *sql.DB, path string, title string, lastModded int64) error {
	_, err := db.Exec(`
		UPDATE books
		SET title = ?, last_modded = ?
		WHERE path = ?
	`, title, lastModded, path)
	return err
}

func UpdateBookPositionAndProgress(db *sql.DB, path string, currentPosition string, progress float64) error {
	_, err := db.Exec(`
		UPDATE books
		SET current_position = ?, progress = ?
		WHERE path = ?
	`, currentPosition, progress, path)
	return err
}

func UpdateBookLastOpened(db *sql.DB, path string) error {
	now := time.Now().Unix()
	_, err := db.Exec(`
		UPDATE books
		SET last_opened = ?
		WHERE path = ?
	`, now, path)
	return err
}

func DeleteBookByPath(db *sql.DB, path string) error {
	_, err := db.Exec(`DELETE FROM books WHERE path = ?`, path)
	return err
}

func GetBooksFlat(db *sql.DB, sortBy, order string, limit, offset int) ([]BookData, error) {
	query := fmt.Sprintf(`
		SELECT path, cover_path, type, title, added_time, last_modded, last_opened, current_position, progress
		FROM books
		ORDER BY %s %s
		LIMIT ? OFFSET ?`, sortBy, order)

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []BookData
	for rows.Next() {
		var b BookData
		if err := rows.Scan(
			&b.Path,
			&b.CoverPath,
			&b.Type,
			&b.Title,
			&b.AddedTime,
			&b.LastModded,
			&b.LastOpened,
			&b.CurrentPosition,
			&b.Progress,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func GetBooksFolder(db *sql.DB, folderPath, sortBy, order string, limit, offset int) ([]BookData, error) {
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	query := fmt.Sprintf(`
		SELECT path, cover_path, type, title, added_time, last_modded, last_opened, current_position, progress
		FROM books
		WHERE path LIKE ? AND
		      LENGTH(REPLACE(SUBSTR(path, LENGTH(?) + 1), '/', '')) = LENGTH(SUBSTR(path, LENGTH(?) + 1))
		ORDER BY %s %s
		LIMIT ? OFFSET ?`, sortBy, order)

	rows, err := db.Query(query, folderPath+"%", folderPath, folderPath, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []BookData
	for rows.Next() {
		var b BookData
		if err := rows.Scan(
			&b.Path,
			&b.CoverPath,
			&b.Type,
			&b.Title,
			&b.AddedTime,
			&b.LastModded,
			&b.LastOpened,
			&b.CurrentPosition,
			&b.Progress,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

type ChildFolder struct {
	Path      string
	CoverPath string
}

func GetSubfolders(db *sql.DB, folderPath string) ([]ChildFolder, error) {
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	rows, err := db.Query(`
		SELECT path, cover_path FROM books
		WHERE path LIKE ?
		ORDER BY path;
	`, folderPath+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type item struct {
		path      string
		coverPath string
	}

	var entries []item
	for rows.Next() {
		var i item
		if err := rows.Scan(&i.path, &i.coverPath); err != nil {
			return nil, err
		}
		entries = append(entries, i)
	}

	groupMap := make(map[string]string)

	for _, e := range entries {
		trimmed := strings.TrimPrefix(e.path, folderPath)
		parts := strings.Split(trimmed, "/")
		if len(parts) < 2 {
			continue
		}
		subfolder := folderPath + parts[0]

		if _, exists := groupMap[subfolder]; !exists {
			groupMap[subfolder] = e.coverPath
		}
	}

	var results []ChildFolder

	var keys []string
	for k := range groupMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		results = append(results, ChildFolder{
			Path:      k,
			CoverPath: groupMap[k],
		})
	}

	return results, nil
}

func SearchBooksInPathListWithTitlePrefix(
	db *sql.DB,
	pathList []string,
	titleLike string,
	sortBy, order string,
	limit, offset int,
) ([]BookData, error) {
	if len(pathList) == 0 {
		return []BookData{}, nil
	}

	placeholders := strings.Repeat("?,", len(pathList))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]interface{}, len(pathList)+3)
	for i, p := range pathList {
		args[i] = p
	}
	args[len(pathList)] = titleLike + "%"
	args[len(pathList)+1] = limit
	args[len(pathList)+2] = offset

	query := fmt.Sprintf(`
		SELECT path, cover_path, type, title, added_time, last_modded, last_opened, current_position, progress
		FROM books
		WHERE path IN (%s) AND title LIKE ?
		ORDER BY %s %s
		LIMIT ? OFFSET ?;
	`, placeholders, sortBy, strings.ToUpper(order))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []BookData
	for rows.Next() {
		var b BookData
		if err := rows.Scan(
			&b.Path,
			&b.CoverPath,
			&b.Type,
			&b.Title,
			&b.AddedTime,
			&b.LastModded,
			&b.LastOpened,
			&b.CurrentPosition,
			&b.Progress,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}

func SearchBooksInPathList(
	db *sql.DB,
	pathList []string,
	sortBy, order string,
	limit, offset int,
) ([]BookData, error) {
	if len(pathList) == 0 {
		return []BookData{}, nil
	}

	placeholders := strings.Repeat("?,", len(pathList))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]interface{}, len(pathList)+2)
	for i, p := range pathList {
		args[i] = p
	}
	args[len(pathList)] = limit
	args[len(pathList)+1] = offset

	query := fmt.Sprintf(`
		SELECT path, cover_path, type, title, added_time, last_modded, last_opened, current_position, progress
		FROM books
		WHERE path IN (%s)
		ORDER BY %s %s
		LIMIT ? OFFSET ?;
	`, placeholders, sortBy, strings.ToUpper(order))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []BookData
	for rows.Next() {
		var b BookData
		if err := rows.Scan(
			&b.Path,
			&b.CoverPath,
			&b.Type,
			&b.Title,
			&b.AddedTime,
			&b.LastModded,
			&b.LastOpened,
			&b.CurrentPosition,
			&b.Progress,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}

func SearchBooksWithTitlePrefix(
	db *sql.DB,
	titleLike string,
	sortBy, order string,
	limit, offset int,
) ([]BookData, error) {
	args := []interface{}{titleLike + "%", limit, offset}

	query := fmt.Sprintf(`
		SELECT path, cover_path, type, title, added_time, last_modded, last_opened, current_position, progress
		FROM books
		WHERE title LIKE ?
		ORDER BY %s %s
		LIMIT ? OFFSET ?;
	`, sortBy, strings.ToUpper(order))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []BookData
	for rows.Next() {
		var b BookData
		if err := rows.Scan(
			&b.Path,
			&b.CoverPath,
			&b.Type,
			&b.Title,
			&b.AddedTime,
			&b.LastModded,
			&b.LastOpened,
			&b.CurrentPosition,
			&b.Progress,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}
