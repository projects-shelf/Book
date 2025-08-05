package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var dpi int

func init() {
	dpi = getEnvInt("PDF_RENDERING_DPI", 300)
}

func getEnvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}
	return def
}

func PDFPagesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathParam := c.Query("path")
		if pathParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'path' parameter"})
			return
		}

		decodedPath, err := url.PathUnescape(pathParam)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		rawPath := strings.TrimPrefix(decodedPath, "/")
		filePath := filepath.Join("/books", rawPath)

		// Check if file exists (consider adding path validation as needed)
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			}
			return
		}

		// Execute the 'pdfinfo' command from poppler utils
		cmd := exec.Command("pdfinfo", filePath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("pdfinfo error: %v", err)})
			return
		}

		// Parse output to find the page count line
		scanner := bufio.NewScanner(bytes.NewReader(out))
		var pageCount int
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Pages:") {
				fmt.Sscanf(line, "Pages: %d", &pageCount)
				break
			}
		}

		// If page count couldn't be parsed, return error
		if pageCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse page count"})
			return
		}

		// Return JSON response with page count
		c.JSON(http.StatusOK, gin.H{
			"pages": pageCount,
		})
	}
}

func PDFStreamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathParam := c.DefaultQuery("path", "")
		pageParam := c.DefaultQuery("page", "")
		if pathParam == "" || pageParam == "" {
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

		page, err := strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			c.String(http.StatusBadRequest, "invalid page number")
			return
		}

		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				c.Status(http.StatusNotFound)
			} else {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		tmpFile := fmt.Sprintf("pdf-%d.tmp.png", page)

		cmdList := exec.Command("gs",
			"-sDEVICE=png16m",
			"-dUseCIEColor=false",
			fmt.Sprintf("-dFirstPage=%d", page),
			fmt.Sprintf("-dLastPage=%d", page),
			fmt.Sprintf("-r%d", dpi),
			"-dNOPAUSE",
			"-dBATCH",
			"-sOutputFile="+tmpFile,
			filePath,
		)

		if err := cmdList.Run(); err != nil {
			log.Printf("ghostscript command failed: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "image/png")
		c.Header("Cache-Control", "public, max-age=3600")

		c.File(tmpFile)

		go func(file string) {
			time.Sleep(2 * time.Second)
			os.Remove(file)
		}(tmpFile)
	}
}
