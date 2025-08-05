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

// extractCBZCover extracts the first image file from a CBZ archive using
// 7z command without full extraction, decodes it, resizes if needed, and saves as WebP.
func extractCBZCover(cbzPath, outputWebPPath string) error {
	outDir := filepath.Dir(outputWebPPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// 1. List archive contents with '7z l -ba'
	cmdList := exec.Command("7z", "l", "-ba", cbzPath)
	outList, err := cmdList.Output()
	if err != nil {
		return fmt.Errorf("7z list command failed: %w", err)
	}

	// Parse the output line-by-line and collect image files
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
		return fmt.Errorf("no image files found in CBZ archive")
	}

	// Sort filenames and select the first image as cover
	sort.Strings(files)
	targetFile := files[0]

	// 2. Extract the selected image file to stdout using '7z x -so'
	cmdExtract := exec.Command("7z", "x", "-so", cbzPath, targetFile)
	stdout, err := cmdExtract.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	if err := cmdExtract.Start(); err != nil {
		return fmt.Errorf("failed to start 7z extract: %w", err)
	}

	// 3. Read the extracted image data from stdout
	imgData, err := io.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("failed to read extracted image data: %w", err)
	}
	if err := cmdExtract.Wait(); err != nil {
		return fmt.Errorf("7z extract command failed: %w", err)
	}

	// 4. Decode image data
	img, err := decodeImageWithWebP(imgData)
	if err != nil {
		return fmt.Errorf("failed to decode image data: %w", err)
	}

	// 5. Resize image if short side is greater than 512 px

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

	// 6. Encode resized image as WebP and save to output file
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
