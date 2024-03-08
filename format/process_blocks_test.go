package format

import (
	"fmt"
	"os"
	"testing"

	"github.com/s-kngstn/notionsync/api" // adjust the import path based on your project structure
	// other imports as needed
)

// MockNotionAPI defined here
type MockNotionAPI struct {
	BlockTitleResponse    string
	FetchBlockTitleError  error
	ChildBlocksResponse   *api.ResultsWrapper
	FetchChildBlocksError error
}

func (m *MockNotionAPI) GetNotionBlockTitle(pageID, bearerToken string) (string, error) {
	// Mock implementation...
	return m.BlockTitleResponse, m.FetchBlockTitleError
}

func (m *MockNotionAPI) GetNotionChildBlocks(blockID, bearerToken string) (*api.ResultsWrapper, error) {
	// Mock implementation...
	return m.ChildBlocksResponse, m.FetchChildBlocksError
}

func TestProcessBlocksMarkdownOutput(t *testing.T) {
	// Setup Mock API with a response
	mockAPI := &MockNotionAPI{
		BlockTitleResponse: "Example Title",
		ChildBlocksResponse: &api.ResultsWrapper{
			Results: []api.Block{},
		},
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "process_blocks_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	outputPath := tempFile.Name()
	defer os.Remove(outputPath) // Clean up after the test

	uuid := "test-uuid"
	bearerToken := "test-token"
	pageName := "Test Page"
	processedBlocks := make(map[string]string)
	results := &api.ResultsWrapper{
		Results: []api.Block{
			// Define a simple block structure for testing
		},
	}

	ProcessBlocks(uuid, results, outputPath, pageName, mockAPI, bearerToken, processedBlocks)

	// Read the content of the temporary file
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	// Check if the content of the file matches expected output
	expectedContent := "# Test Page\n\n" // Adjust this based on your expected Markdown output
	if string(content) != expectedContent {
		t.Errorf("Expected file content to be %v, got %v", expectedContent, string(content))
	}
}

func TestProcessBlocksWithFetchBlockTitleError(t *testing.T) {
	// Setup Mock API to return an error for FetchBlockTitle
	mockAPI := &MockNotionAPI{
		FetchBlockTitleError: fmt.Errorf("simulated error"),
	}

	// Other setup steps similar to the first test
	// Create a temporary file for output
	tempFile, err := os.CreateTemp("", "process_blocks_error_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	outputPath := tempFile.Name()
	defer os.Remove(outputPath) // Clean up after the test

	uuid := "error-test-uuid"
	bearerToken := "error-test-token"
	pageName := "Error Test Page"
	processedBlocks := make(map[string]string)
	results := &api.ResultsWrapper{
		Results: []api.Block{
			{ID: "error1", Type: "link_to_page", LinkToPage: &api.LinkToPage{PageID: "errorPage"}},
		},
	}

	ProcessBlocks(uuid, results, outputPath, pageName, mockAPI, bearerToken, processedBlocks)

	// Error handling test logic remains the same
}
