package api

import (
	"back/database"
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RootHandler(bookDB *sql.DB, pageSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		pathParam := c.Param("path")
		folderPath := filepath.Join("/books", pathParam)

		sortBy := c.DefaultQuery("sort", "title")
		order := c.DefaultQuery("order", "asc")
		pageStr := c.DefaultQuery("page", "1")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		validSorts := map[string]string{
			"title":       "title",
			"added_time":  "added_time",
			"last_opened": "last_opened",
			"progress":    "progress",
		}

		sortColumn, ok := validSorts[sortBy]
		if !ok {
			sortColumn = "title"
		}

		order = strings.ToUpper(order)
		if order != "ASC" && order != "DESC" {
			order = "ASC"
		}

		// folder

		var subfolders []database.ChildFolder
		if page == 1 {
			subfolders, err = database.GetSubfolders(bookDB, folderPath)
			if err != nil {
				fmt.Println("failed to get subfolders: %v", err)
			}
		}

		// book

		offset := (page - 1) * pageSize

		booksData, err := database.GetBooksFolder(bookDB, folderPath, sortColumn, order, pageSize+1, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
			return
		}

		hasMore := false
		if len(booksData) > pageSize {
			hasMore = true
			booksData = booksData[:pageSize]
		}

		subNum := len(subfolders)
		books := make([]BookEntry, subNum+len(booksData))
		for i, f := range subfolders {
			books[i] = BookEntry{
				Type:            "Folder",
				Path:            strings.TrimPrefix(f.Path, "/books"),
				Cover:           strings.TrimPrefix(f.CoverPath, "/cache/covers"),
				Title:           path.Base(f.Path),
				CurrentPosition: "",
				Progress:        0.0,
			}
		}

		for i, b := range booksData {
			books[i+subNum] = BookEntry{
				Type:            b.Type,
				Path:            strings.TrimPrefix(b.Path, "/books"),
				Cover:           strings.TrimPrefix(b.CoverPath, "/cache/covers"),
				Title:           b.Title,
				CurrentPosition: b.CurrentPosition,
				Progress:        b.Progress,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"books":   books,
			"hasMore": hasMore,
		})
	}
}
