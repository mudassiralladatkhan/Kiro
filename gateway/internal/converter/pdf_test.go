package converter

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestConvertAnthropicContentToText_DocumentPlain(t *testing.T) {
	plainText := "Hello, this is a plain text file upload test!"
	base64Data := base64.StdEncoding.EncodeToString([]byte(plainText))

	content := []any{
		map[string]any{
			"type": "document",
			"name": "test.txt",
			"source": map[string]any{
				"type":       "base64",
				"media_type": "text/plain",
				"data":       base64Data,
			},
		},
	}

	extracted := convertAnthropicContentToText(content)
	if !strings.Contains(extracted, "Uploaded Document Content: test.txt") {
		t.Fatalf("expected document title in extracted text, got: %s", extracted)
	}
	if !strings.Contains(extracted, plainText) {
		t.Fatalf("expected document content in extracted text, got: %s", extracted)
	}
}

func TestConvertAnthropicContentToText_Unsupported(t *testing.T) {
	content := []any{
		map[string]any{
			"type": "document",
			"name": "test.docx",
			"source": map[string]any{
				"type":       "base64",
				"media_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"data":       "abcdef",
			},
		},
	}

	extracted := convertAnthropicContentToText(content)
	if !strings.Contains(extracted, "Unsupported document type") {
		t.Fatalf("expected unsupported type warning in extracted text, got: %s", extracted)
	}
}
