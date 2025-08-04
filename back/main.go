package main

import (
	"back/api"
	"back/api/stream"
	"back/database"
	"back/internal/scan"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := os.MkdirAll("/db/", 0755); err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// db

	bookDB := database.OpenBookDB()
	defer bookDB.Close()

	keywordDB := database.OpenKeywordDB()
	defer keywordDB.Close()

	// scan books

	if err := scan.Scan(bookDB, keywordDB); err != nil {
		panic(err)
	}

	// api

	pageSize := 20
	if v := os.Getenv("PAGE_SIZE"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	r.GET("/api/all", api.AllHandler(bookDB, pageSize))
	r.GET("/api/root/*path", api.RootHandler(bookDB, pageSize))
	r.GET("/api/search", api.SearchHandler(bookDB, keywordDB, pageSize))
	r.GET("/api/progress", api.ProgressHandler(bookDB))
	r.GET("/api/access", api.AccessHandler(bookDB))

	r.GET("/book/pdf", stream.PDFStreamHandler())
	r.GET("/book/epub", stream.EPUBStreamHandler())

	r.GET("/book/cbr", stream.CBRStreamHandler())
	r.GET("/book/cbr/pages", stream.CBRPagesHandler())
	r.GET("/book/cbz", stream.CBZStreamHandler())
	r.GET("/book/cbz/pages", stream.CBZPagesHandler())

	r.GET("/cover/*path", api.CoverHandler())

	r.Run(":8080")
}
