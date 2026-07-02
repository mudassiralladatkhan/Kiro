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

func TestContainsVisualKeywords(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
	}{
		{"Look at the style of this PDF", true},
		{"Tell me if the design looks good", true},
		{"What does this say?", false},
		{[]any{map[string]any{"type": "text", "text": "Analyze the visual elements"}}, true},
		{[]any{map[string]any{"type": "text", "text": "Just summarize the text"}}, false},
	}
	for _, tt := range tests {
		result := containsVisualKeywords(tt.input)
		if result != tt.expected {
			t.Errorf("containsVisualKeywords(%v) = %v; expected %v", tt.input, result, tt.expected)
		}
	}
}
