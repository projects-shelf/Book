package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func CoverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathParam := c.Param("path")
		if pathParam == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		cleanPath := pathParam[1:]

		filePath := filepath.Join("/cache", "covers", cleanPath)

		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				c.Status(http.StatusNotFound)
			} else {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		c.Header("Cache-Control", "public, max-age=3600")
		c.File(filePath)
	}
}
