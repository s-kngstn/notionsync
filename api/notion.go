package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Implement GetRichText for Heading
func (h *Heading) GetRichText() []RichText {
	return h.RichText
}

// Implement GetRichText for Paragraph
func (p *Paragraph) GetRichText() []RichText {
	return p.RichText
}

// HttpClientInterface defines the interface for the HTTP client
type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// NotionApiClient struct will hold any dependencies for your API client, e.g., the HTTP client
type NotionApiClient struct {
	Client HttpClientInterface
}

// NewNotionApiClient creates a new API client with the provided HTTP client
func NewNotionApiClient(client HttpClientInterface) *NotionApiClient {
	return &NotionApiClient{
		Client: client,
	}
}

// CallAPI performs the actual API call and processes the response
func (api *NotionApiClient) CallAPI(customID, bearerToken string) error {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", customID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// @todo we'll need the user to provide their own bearer token
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Notion-Version", "2022-06-28")

	resp, err := api.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		var apiError APIErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("error parsing API error response: %w", err)
		}
		return fmt.Errorf("API Error: %s - %s", apiError.Code, apiError.Message)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var results ResultsWrapper
	if err := json.Unmarshal(body, &results); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Open a new file for writing, create if not exists, truncate if exists
	outputPath := "test.md"
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating markdown file: %w", err)
	}
	defer file.Close()

	for _, block := range results.Results {
		fmt.Printf("Block ID: %s\n", block.ID)
		fmt.Printf("Block Type: %s\n", block.Type)

		var provider RichTextProvider

		switch block.Type {
		case "heading_1":
			provider = block.Heading1
		case "heading_2":
			provider = block.Heading2
		case "heading_3":
			provider = block.Heading3
		case "paragraph":
			provider = block.Paragraph
		}

		if provider != nil {
			for _, rt := range provider.GetRichText() {
				// @TODO: Depending on the type of Rich Text, I need to format the content differently
				_, err := file.WriteString(rt.Text.Content + "\n")
				if err != nil {
					return fmt.Errorf("error writing to markdown file: %w", err)
				}
				fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
			}
		}
	}

	return nil // Return nil if everything was successful
}
