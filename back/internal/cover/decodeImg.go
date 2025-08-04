package cover

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"

	"github.com/chai2010/webp"
)

func decodeImageWithWebP(imgData []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err == nil {
		return img, nil
	}

	// WebP fallback
	img, err = webp.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (even with WebP): %w", err)
	}

	return img, nil
}
