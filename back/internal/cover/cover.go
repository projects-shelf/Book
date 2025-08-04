package cover

import "fmt"

const size = 200
const quality = 70

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
