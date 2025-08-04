package cover

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

func extractPDFCover(pdfPath, outputWebPPath string) error {
	outDir := filepath.Dir(outputWebPPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	tmpJPEG := outputWebPPath + ".tmp.jpg"

	// Step 1: extract first page as JPEG
	// pdftoppm expects a prefix, not a filename
	prefix := strings.TrimSuffix(tmpJPEG, ".jpg")
	cmd := exec.Command("pdftoppm", "-jpeg", "-f", "1", "-l", "1", "-r", "150", pdfPath, prefix)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pdftoppm failed: %v\noutput: %s", err, out)
	}

	// Output file will be prefix + "-[0]*1.jpg"
	matches, err := filepath.Glob(prefix + "-*.jpg")
	if err != nil || len(matches) == 0 {
		return fmt.Errorf("no jpg files found matching pattern: %v", prefix+"-*.jpg")
	}
	jpgPath := matches[0]

	// Step 2: open and decode the image
	f, err := os.Open(jpgPath)
	if err != nil {
		return fmt.Errorf("failed to open JPEG: %w", err)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode JPEG: %w", err)
	}

	// Step 3: resize if needed (short side ≥ Size)
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

	// Step 4: save as WebP
	outFile, err := os.Create(outputWebPPath)
	if err != nil {
		return fmt.Errorf("failed to create WebP file: %w", err)
	}
	defer outFile.Close()

	var buf bytes.Buffer
	if err := webp.Encode(&buf, resized, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return fmt.Errorf("failed to encode WebP: %w", err)
	}
	if _, err := buf.WriteTo(outFile); err != nil {
		return fmt.Errorf("failed to write WebP file: %w", err)
	}

	// Step 5: cleanup
	_ = os.Remove(jpgPath)

	return nil
}
