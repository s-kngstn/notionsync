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
			got := toTitleCase(tt.input)
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
		{
			name: "Strikethrough Annotation",
			rt: api.RichText{
				Text: api.Text{
					Content: "Test",
				},
				Annotations: api.Annotations{
					Strikethrough: true,
				},
			},
			expected: "~~Test~~",
		},
		{
			name: "Code Annotation",
			rt: api.RichText{
				Text: api.Text{
					Content: "Test",
				},
				Annotations: api.Annotations{
					Code: true,
				},
			},
			expected: "`Test`",
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
			{
				ID:   "8",
				Type: "to_do",
				Todo: &api.Todo{
					RichText: []api.RichText{{Text: api.Text{Content: "To Do Item"}}},
					Checked:  true,
				},
			},
			{
				ID:   "8",
				Type: "to_do",
				Todo: &api.Todo{
					RichText: []api.RichText{{Text: api.Text{Content: "To Do Item unchecked"}}},
					Checked:  false,
				},
			},
			{
				ID:   "9",
				Type: "code",
				Code: &api.Code{
					RichText: []api.RichText{{Text: api.Text{Content: "Code Block"}}},
					Language: "go",
				},
			},
			{
				ID:   "10",
				Type: "bookmark",
				Bookmark: &api.Bookmark{
					URL: "https://example.com",
				},
			},
			{
				ID:   "11",
				Type: "quote",
				Quote: &api.Quote{
					RichText: []api.RichText{{Text: api.Text{Content: "Quote Block"}}},
				},
			},
			{
				ID:      "12",
				Type:    "divider",
				Divider: &api.Divider{},
			},
		},
		// your block definitions...
	}

	outputPath := "./test_output.md"
	pageName := "Test Page"
	pageTitles := map[string]string{}
	if err := WriteBlocksToMarkdown(results, outputPath, pageName, pageTitles); err != nil {
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

	// Expected content checks
	expectedContents := []string{
		"# Test Page\n\n",
		"Hello World\n",
		"# Heading One\n",
		"## Heading Two\n",
		"### Heading Three\n",
		"- List Item One\n",
		"1. List Item Two\n",
		"2. List Item Three\n",
		"- [x] To Do Item\n",
		"- [ ] To Do Item unchecked\n",
		"```go\nCode Block\n ```\n",
		"- [https://example.com]\n",
		"> Quote Block\n",
		"---\n",
	}

	for _, ec := range expectedContents {
		if !strings.Contains(string(content), ec) {
			t.Errorf("File content does not contain expected text: %q", ec)
		}
	}

	// Clean up
	os.Remove(outputPath)
}
