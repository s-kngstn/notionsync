package format

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/s-kngstn/notionsync/api"
)

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", "Test"},
		{"multiple words test", "Multiple Words Test"},
		{"test-hyphenated-words", "Test Hyphenated Words"},
		{"mIxEd CaSe", "Mixed Case"},
		{"testing with in and on", "Testing With In And On"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toTitleCase(tt.input) // Call your toTitleCase function
			if got != tt.expected {
				t.Errorf("toTitleCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

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
	pageName := "Test Page" // Add a mock page name for testing
	if err := WriteBlocksToMarkdown(results, outputPath, pageName); err != nil {
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

	// Check if the content contains the expected page title and body text
	expectedTitle := "# Test Page\n\n" // Adjust according to how the title is formatted
	if !strings.Contains(string(content), expectedTitle) || !strings.Contains(string(content), "Hello World") {
		t.Errorf("File content does not contain expected text. Expected title %q and body text 'Hello World'. Got: %v", expectedTitle, string(content))
	}

	// Clean up
	os.Remove(outputPath)
}
