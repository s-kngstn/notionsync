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
	results := &api.ResultsWrapper{
		Results: []api.Block{
			{
				ID:   "1",
				Type: "paragraph",
				Paragraph: &api.Paragraph{
					RichText: []api.RichText{{Text: api.Text{Content: "Hello World"}}},
				},
			},
			{
				ID:   "2",
				Type: "heading_1",
				Heading1: &api.Heading{
					RichText: []api.RichText{{Text: api.Text{Content: "Heading One"}}},
				},
			},
			{
				ID:   "3",
				Type: "heading_2",
				Heading2: &api.Heading{
					RichText: []api.RichText{{Text: api.Text{Content: "Heading Two"}}},
				},
			},
			{
				ID:   "4",
				Type: "heading_3",
				Heading3: &api.Heading{
					RichText: []api.RichText{{Text: api.Text{Content: "Heading Three"}}},
				},
			},
			{
				ID:   "5",
				Type: "bulleted_list_item",
				Bulleted: &api.ListItem{
					RichText: []api.RichText{{Text: api.Text{Content: "List Item One"}}},
				},
			},
			{
				ID:   "6",
				Type: "numbered_list_item",
				Numbered: &api.ListItem{
					RichText: []api.RichText{{Text: api.Text{Content: "List Item Two"}}},
				},
			},
			{
				ID:   "7",
				Type: "numbered_list_item",
				Numbered: &api.ListItem{
					RichText: []api.RichText{{Text: api.Text{Content: "List Item Three"}}},
				},
			},
		},
	}

	outputPath := "./test_output.md"
	pageName := "Test Page"
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

	// Check for the expected title, paragraph, and all headings and list item in the content
	expectedTitle := "# Test Page\n\n"
	expectedHeading1 := "# Heading One\n"
	expectedHeading2 := "## Heading Two\n"
	expectedHeading3 := "### Heading Three\n"
	expectedListItem := "- List Item One\n"
	expectedNumberedListItem := "1. List Item Two\n"
	expectedNumberedListItem2 := "2. List Item Three\n"
	if !strings.Contains(string(content), expectedTitle) ||
		!strings.Contains(string(content), "Hello World") ||
		!strings.Contains(string(content), expectedHeading1) ||
		!strings.Contains(string(content), expectedHeading2) ||
		!strings.Contains(string(content), expectedHeading3) ||
		!strings.Contains(string(content), expectedListItem) ||
		!strings.Contains(string(content), expectedNumberedListItem) ||
		!strings.Contains(string(content), expectedNumberedListItem2) {
		t.Errorf("File content does not contain expected text. Expected title %q, body text 'Hello World', headings %q, %q, %q, list items %q, %q, %q Got: %v",
			expectedTitle, expectedHeading1, expectedHeading2, expectedHeading3, expectedListItem, expectedNumberedListItem, expectedNumberedListItem2, string(content))
	}

	// Clean up
	os.Remove(outputPath)
}
