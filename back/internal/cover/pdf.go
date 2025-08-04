package cover

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

func extractPDFCover(pdfPath, outputWebPPath string) error {
	outDir := filepath.Dir(outputWebPPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	tmpPNG := outputWebPPath + ".tmp.png"

	// Step 1: Convert the first page of the PDF to PNG using Ghostscript
	cmd := exec.Command("gs",
		"-sDEVICE=png16m",
		"-dFirstPage=1",
		"-dLastPage=1",
		"-r150",
		"-o", tmpPNG,
		pdfPath,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ghostscript failed: %v\noutput: %s", err, out)
	}

	// Step 2: open and decode the PNG image
	f, err := os.Open(tmpPNG)
	if err != nil {
		return fmt.Errorf("failed to open PNG: %w", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Step 3: resize if needed (short side ≥ size)
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
	_ = os.Remove(tmpPNG)

	return nil
}
