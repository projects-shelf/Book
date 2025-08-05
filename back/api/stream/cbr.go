package stream

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CBRPagesHandler() gin.HandlerFunc {
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

		cmdList := exec.Command("7z", "l", "-ba", filePath)
		outList, err := cmdList.Output()
		if err != nil {
			log.Printf("7z list command failed: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(outList))
		var files []string
		imgExts := []string{".jpg", ".jpeg", ".png", ".webp", ".bmp"}
		for scanner.Scan() {
			line := scanner.Text()
			cols := strings.Fields(line)
			if len(cols) < 6 {
				continue
			}
			name := cols[len(cols)-1]
			ext := strings.ToLower(filepath.Ext(name))
			for _, e := range imgExts {
				if ext == e {
					files = append(files, name)
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("scanner error: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(files) == 0 {
			log.Printf("no image files found in CBR archive")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pages": len(files),
		})
	}
}

func CBRStreamHandler() gin.HandlerFunc {
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
		if err != nil {
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

		cmdList := exec.Command("7z", "l", "-ba", filePath)
		outList, err := cmdList.Output()
		if err != nil {
			log.Printf("7z list command failed: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(outList))
		var files []string
		imgExts := []string{".jpg", ".jpeg", ".png", ".webp", ".bmp"}
		for scanner.Scan() {
			line := scanner.Text()
			cols := strings.Fields(line)
			if len(cols) < 6 {
				continue
			}
			name := cols[len(cols)-1]
			ext := strings.ToLower(filepath.Ext(name))
			for _, e := range imgExts {
				if ext == e {
					files = append(files, name)
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("scanner error: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(files) == 0 {
			log.Printf("no image files found in CBR archive")
			c.Status(http.StatusInternalServerError)
			return
		}

		sort.Strings(files)
		if page < 1 || page > len(files) {
			c.String(http.StatusBadRequest, "invalid page number")
			return
		}
		targetFile := files[page-1]

		cmdExtract := exec.Command("7z", "x", "-so", filePath, targetFile)
		stdout, err := cmdExtract.StdoutPipe()
		if err != nil {
			log.Printf("failed to get stdout pipe: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		if err := cmdExtract.Start(); err != nil {
			log.Printf("failed to start 7z extract: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		mime := detectImageTypeByExt(targetFile)

		c.Header("Content-Type", mime)
		c.Header("Cache-Control", "public, max-age=3600")

		c.Status(http.StatusOK)
		_, copyErr := io.Copy(c.Writer, stdout)

		if copyErr != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := cmdExtract.Wait(); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}

func detectImageTypeByExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	default:
		return "application/octet-stream"
	}
}
