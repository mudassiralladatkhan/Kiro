package converter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/rs/zerolog/log"
)

// extractTextFromPDFBase64 decodes a base64-encoded PDF, parses it, and
// extracts all of its plain text.
func extractTextFromPDFBase64(base64Data string) (string, error) {
	// 1. Decode base64 data
	pdfBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	// 2. Initialize bytes reader (ReaderAt)
	reader := bytes.NewReader(pdfBytes)
	size := int64(len(pdfBytes))

	// 3. Parse PDF structure
	pdfReader, err := pdf.NewReader(reader, size)
	if err != nil {
		return "", fmt.Errorf("parse pdf: %w", err)
	}

	// 4. Extract plain text
	textReader, err := pdfReader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("get plain text reader: %w", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, textReader); err != nil {
		// Sometimes GetPlainText can fail mid-read on malformed PDFs,
		// but we still want to return whatever text we successfully read.
		log.Warn().Err(err).Msg("partial error during PDF text extraction")
	}

	return buf.String(), nil
}

// renderPDFToPNGs decodes a base64-encoded PDF, renders its pages to PNG files
// using pdftoppm, and returns the pages as base64-encoded PNG image data.
// It limits the pages to maxPages to prevent hitting payload/token limits.
func renderPDFToPNGs(base64Data string, maxPages int) ([]string, error) {
	// Check if pdftoppm is available on the system path
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return nil, fmt.Errorf("pdftoppm executable not found in path: %w", err)
	}

	// 1. Decode base64 data
	pdfBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	// 2. Create a temporary file for the PDF
	tmpFile, err := os.CreateTemp("", "kiro-pdf-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(pdfBytes); err != nil {
		return nil, fmt.Errorf("write temp pdf: %w", err)
	}
	tmpFile.Close()

	// 3. Create a temporary directory for output PNG files
	tmpDir, err := os.MkdirTemp("", "kiro-pdf-pages-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 4. Run pdftoppm
	// -png: output PNG format
	// -r 150: 150 DPI resolution (good balance between quality and payload size)
	// -l maxPages: limit to maxPages pages
	outputPrefix := filepath.Join(tmpDir, "page")
	cmd := exec.Command("pdftoppm", "-png", "-r", "150", "-l", fmt.Sprintf("%d", maxPages), tmpFile.Name(), outputPrefix)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run pdftoppm: %w (stderr: %s)", err, stderr.String())
	}

	// 5. Read output PNG files in sorted order
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("read temp dir: %w", err)
	}

	var pngFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".png") {
			pngFiles = append(pngFiles, filepath.Join(tmpDir, f.Name()))
		}
	}
	sort.Strings(pngFiles)

	var base64PNGs []string
	for _, path := range pngFiles {
		imgBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read page file %s: %w", path, err)
		}
		base64PNGs = append(base64PNGs, base64.StdEncoding.EncodeToString(imgBytes))
	}

	return base64PNGs, nil
}
