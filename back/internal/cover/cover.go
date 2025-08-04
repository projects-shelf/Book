package cover

import (
	"fmt"
	"os"
	"strconv"
)

var (
	size    int
	quality int
)

func init() {
	size = getEnvInt("COVER_SIZE", 200)
	quality = getEnvInt("COVER_QUALITY", 70)
}

func getEnvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}
	return def
}

func ExtractCover(inputPath, outputPath, bookType string) error {
	switch bookType {
	case "PDF":
		return extractPDFCover(inputPath, outputPath)
	case "EPUB":
		return extractEPUBCover(inputPath, outputPath)
	case "CBZ":
		return extractCBZCover(inputPath, outputPath)
	case "CBR":
		return extractCBRCover(inputPath, outputPath)
	default:
		return fmt.Errorf("unsupported file type: %s", bookType)
	}
}
