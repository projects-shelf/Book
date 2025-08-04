package meta

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type EPUBMeta struct {
	Title    string
	Keywords []string
	ModTime  time.Time
}

// ExtractEPUBMeta extracts metadata and last modified time from the given EPUB file
func extractEPUBMeta(path string) (string, []string, int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", nil, 0, fmt.Errorf("failed to stat file: %w", err)
	}
	modTime := info.ModTime().Unix()

	// Step 1: Get path to content.opf
	containerXML, err := run7zCommand(path, "META-INF/container.xml")
	if err != nil {
		return "", nil, 0, fmt.Errorf("failed to extract container.xml: %w", err)
	}

	var container struct {
		Rootfiles struct {
			Rootfile struct {
				FullPath string `xml:"full-path,attr"`
			} `xml:"rootfile"`
		} `xml:"rootfiles"`
	}

	if err := xml.Unmarshal(containerXML, &container); err != nil {
		return "", nil, 0, fmt.Errorf("failed to parse container.xml: %w", err)
	}

	contentPath := container.Rootfiles.Rootfile.FullPath
	if contentPath == "" {
		return "", nil, 0, fmt.Errorf("content.opf path not found in container.xml")
	}

	// Step 2: Extract content.opf
	opfData, err := run7zCommand(path, contentPath)
	if err != nil {
		return "", nil, 0, fmt.Errorf("failed to extract content.opf: %w", err)
	}

	// Step 3: Parse metadata from content.opf
	type metadata struct {
		Title   string   `xml:"title"`
		Creator string   `xml:"creator"`
		Subject []string `xml:"subject"`
	}

	var pkg struct {
		Metadata metadata `xml:"metadata"`
	}

	if err := xml.Unmarshal(opfData, &pkg); err != nil {
		return "", nil, 0, fmt.Errorf("failed to parse content.opf: %w", err)
	}

	title := pkg.Metadata.Title
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	var keywords []string
	if pkg.Metadata.Creator != "" {
		keywords = append(keywords, pkg.Metadata.Creator)
	}
	if len(pkg.Metadata.Subject) > 0 {
		keywords = append(keywords, pkg.Metadata.Subject...)
	}

	return title, keywords, modTime, nil
}

// run7zCommand extracts a single file from the EPUB archive using 7z and returns its contents.
func run7zCommand(epubPath, internalPath string) ([]byte, error) {
	cmd := exec.Command("7z", "x", "-so", epubPath, internalPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
