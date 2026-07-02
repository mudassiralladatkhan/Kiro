package converter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

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
