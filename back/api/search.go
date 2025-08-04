package api

import (
	"back/database"
	"database/sql"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SearchHandler(bookDB, keywordDB *sql.DB, pageSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			"last_opened": "last_access_time",
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

		// query

		q := c.DefaultQuery("q", "")
		decodedQ, err := url.QueryUnescape(q)
		if err != nil {
			decodedQ = q
		}

		parts := strings.Split(decodedQ, ",")

		var titleQuery string
		var keywordsList []string

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			if strings.HasPrefix(part, "#") {
				keywordsList = append(keywordsList, strings.TrimPrefix(part, "#"))
			} else if titleQuery == "" {
				titleQuery = part
			}
		}

		offset := (page - 1) * pageSize

		var booksData []database.BookData
		switch {
		case titleQuery != "" && len(keywordsList) != 0:
			pathList, err := database.FindPathsByKeywords(keywordDB, keywordsList)
			if err != nil {
				break
			}
			booksData, err = database.SearchBooksInPathListWithTitlePrefix(bookDB, pathList, titleQuery, sortColumn, order, pageSize+1, offset)
		case titleQuery == "" && len(keywordsList) != 0:
			pathList, err := database.FindPathsByKeywords(keywordDB, keywordsList)
			if err != nil {
				break
			}
			booksData, err = database.SearchBooksInPathList(bookDB, pathList, sortColumn, order, pageSize+1, offset)
		case titleQuery != "" && len(keywordsList) == 0:
			booksData, err = database.SearchBooksWithTitlePrefix(bookDB, titleQuery, sortColumn, order, pageSize+1, offset)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
			return
		}

		hasMore := false
		if len(booksData) > pageSize {
			hasMore = true
			booksData = booksData[:pageSize]
		}

		books := make([]BookEntry, len(booksData))
		for i, b := range booksData {
			books[i] = BookEntry{
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
