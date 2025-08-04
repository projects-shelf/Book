package api

import (
	"back/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccessHandler(bookDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.DefaultQuery("path", "")

		database.UpdateBookLastOpened(bookDB, "/books"+path)

		c.JSON(http.StatusOK, "")
	}
}
