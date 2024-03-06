package format

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/s-kngstn/notionsync/api"
)

func TestApplyAnnotationsToContent(t *testing.T) {
	tests := []struct {
		name     string
		rt       api.RichText
		expected string
	}{
		{
			name: "Bold Annotation",
			rt: api.RichText{
				Text: api.Text{
					Content: "Test",
				},
				Annotations: api.Annotations{
					Bold: true,
				},
			},
			expected: "**Test**",
		},
		{
			name: "Italic Annotation",
			rt: api.RichText{
				Text: api.Text{
					Content: "Test",
				},
				Annotations: api.Annotations{
					Italic: true,
				},
			},
			expected: "*Test*",
		},
		{
			name: "Bold and Italic Annotation",
			rt: api.RichText{
				Text: api.Text{
					Content: "Test",
				},
				Annotations: api.Annotations{
					Bold:   true,
					Italic: true,
				},
			},
			expected: "***Test***",
		},
		// Add more tests for other annotations
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyAnnotationsToContent(tt.rt)
			if result != tt.expected {
				t.Errorf("applyAnnotationsToContent(%v) = %v, want %v", tt.rt, result, tt.expected)
			}
		})
	}
}

func TestWriteBlocksToMarkdown(t *testing.T) {
	// Mock data
	results := &api.ResultsWrapper{
		Results: []api.Block{
			{
				ID:   "1",
				Type: "paragraph",
				Paragraph: &api.Paragraph{
					RichText: []api.RichText{
						{
							Text: api.Text{
								Content: "Hello World",
							},
						},
					},
				},
			},
			// Add more blocks as needed
		},
	}

	outputPath := "./test_output.md"
	if err := WriteBlocksToMarkdown(results, outputPath); err != nil {
		t.Errorf("WriteBlocksToMarkdown returned an error: %v", err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}
	if !strings.Contains(string(content), "Hello World") {
		t.Errorf("File content does not contain expected text. Got: %v", string(content))
	}

	// Clean up
	os.Remove(outputPath)
}
