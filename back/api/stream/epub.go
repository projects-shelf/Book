package stream

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func EPUBStreamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathParam := c.DefaultQuery("path", "")
		if pathParam == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		decodedPath, err := url.PathUnescape(pathParam)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		rawPath := strings.TrimPrefix(decodedPath, "/")
		filePath := filepath.Join("/books", rawPath)

		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				c.Status(http.StatusNotFound)
			} else {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		c.Header("Content-Type", "application/epub+zip")

		c.File(filePath)
	}
}
