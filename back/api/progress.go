package api

import (
	"back/database"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ProgressHandler(bookDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.DefaultQuery("path", "")
		currentPosition := c.DefaultQuery("position", "")
		progressStr := c.DefaultQuery("progress", "0")

		progress, err := strconv.ParseFloat(progressStr, 64)
		if err != nil {
			progress = 0.0
		}

		database.UpdateBookPositionAndProgress(bookDB, "/books"+path, currentPosition, progress)

		c.JSON(http.StatusOK, "")
	}
}
