package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
			fmt.Errorf("7z list command failed: %w", err)
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
			fmt.Errorf("scanner error: %w", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(files) == 0 {
			fmt.Errorf("no image files found in CBR archive")
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
			fmt.Errorf("7z list command failed: %w", err)
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
			fmt.Errorf("scanner error: %w", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(files) == 0 {
			fmt.Errorf("no image files found in CBR archive")
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
			fmt.Errorf("failed to get stdout pipe: %w", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		if err := cmdExtract.Start(); err != nil {
			fmt.Errorf("failed to start 7z extract: %w", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		imgData, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Errorf("failed to read extracted image data: %w", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		if err := cmdExtract.Wait(); err != nil {
			fmt.Errorf("7z extract command failed: %w", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		mime := detectImageTypeByExt(targetFile)

		c.Header("Content-Type", mime)

		c.Data(http.StatusOK, mime, imgData)
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
