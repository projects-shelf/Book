package cover

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

// extractEPUBCover extracts the first image from an EPUB file using 7z and saves it as a WebP.
func extractEPUBCover(epubPath, outputWebPPath string) error {
	outDir := filepath.Dir(outputWebPPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// 1. List archive contents
	cmdList := exec.Command("7z", "l", "-ba", epubPath)
	outList, err := cmdList.Output()
	if err != nil {
		return fmt.Errorf("7z list command failed: %w", err)
	}

	// 2. Parse image files from EPUB (usually in OEBPS/images/ etc.)
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
		return fmt.Errorf("scanner error: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no image files found in EPUB archive")
	}

	sort.Strings(files)
	targetFile := files[0]

	// 3. Extract image via 7z
	cmdExtract := exec.Command("7z", "x", "-so", epubPath, targetFile)
	stdout, err := cmdExtract.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	if err := cmdExtract.Start(); err != nil {
		return fmt.Errorf("failed to start 7z extract: %w", err)
	}
	imgData, err := io.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}
	if err := cmdExtract.Wait(); err != nil {
		return fmt.Errorf("7z extract command failed: %w", err)
	}

	// 4. Decode image
	img, err := decodeImageWithWebP(imgData)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// 5. Resize if needed
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	short := width
	if height < width {
		short = height
	}
	resized := img
	if short > size {
		scale := float64(size) / float64(short)
		newW := int(float64(width) * scale)
		newH := int(float64(height) * scale)
		dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
		draw.ApproxBiLinear.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)
		resized = dst
	}

	// 6. Encode to WebP
	outFile, err := os.Create(outputWebPPath)
	if err != nil {
		return fmt.Errorf("failed to create WebP file: %w", err)
	}
	defer outFile.Close()

	var buf bytes.Buffer
	if err := webp.Encode(&buf, resized, &webp.Options{Lossless: false, Quality: float32(quality)}); err != nil {
		return fmt.Errorf("failed to encode WebP: %w", err)
	}
	if _, err := buf.WriteTo(outFile); err != nil {
		return fmt.Errorf("failed to write WebP file: %w", err)
	}

	return nil
}
