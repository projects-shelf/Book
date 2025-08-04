package meta

import "fmt"

func ExtractMeta(path, bookType string) (string, []string, int64, error) {
	switch bookType {
	case "PDF":
		return extractPDFMeta(path)
	case "EPUB":
		return extractEPUBMeta(path)
	case "CBZ":
		return extractCBZMeta(path)
	case "CBR":
		return extractCBRMeta(path)
	default:
		return "", nil, 0, fmt.Errorf("unsupported file type: %s", bookType)
	}
}
