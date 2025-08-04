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

func extractCBRCover(cbrPath, outputWebPPath string) error {
	outDir := filepath.Dir(outputWebPPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// 1. List archive contents with '7z l -ba'
	cmdList := exec.Command("7z", "l", "-ba", cbrPath)
	outList, err := cmdList.Output()
	if err != nil {
		return fmt.Errorf("7z list command failed: %w", err)
	}

	// 2. Parse image entries
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
		return fmt.Errorf("no image files found in CBR archive")
	}

	// 3. Sort and pick the first image
	sort.Strings(files)
	targetFile := files[0]

	// 4. Extract the selected image file to stdout using '7z x -so'
	cmdExtract := exec.Command("7z", "x", "-so", cbrPath, targetFile)
	stdout, err := cmdExtract.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	if err := cmdExtract.Start(); err != nil {
		return fmt.Errorf("failed to start 7z extract: %w", err)
	}

	imgData, err := io.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("failed to read extracted image data: %w", err)
	}
	if err := cmdExtract.Wait(); err != nil {
		return fmt.Errorf("7z extract command failed: %w", err)
	}

	// 5. Decode image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return fmt.Errorf("failed to decode image data: %w", err)
	}

	// 6. Resize if necessary
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	short := width
	if height < width {
		short = height
	}

	var resized image.Image = img
	if short > size {
		scale := float64(size) / float64(short)
		newW := int(float64(width) * scale)
		newH := int(float64(height) * scale)
		dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
		draw.ApproxBiLinear.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)
		resized = dst
	}

	// 7. Encode to WebP
	outFile, err := os.Create(outputWebPPath)
	if err != nil {
		return fmt.Errorf("failed to create WebP output: %w", err)
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
