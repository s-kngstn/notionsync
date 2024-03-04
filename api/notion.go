package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// APIErrorResponse represents JSON error responses from the API
type APIErrorResponse struct {
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ResultsWrapper is assumed to be the structure of your successful response
type ResultsWrapper struct {
	Results []Block `json:"results"`
}

type Block struct {
	ID        string `json:"id"`
	Paragraph struct {
		RichText []RichText `json:"rich_text"`
	} `json:"paragraph"`
}

type RichText struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
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

	// Process the results as needed
	for _, block := range results.Results {
		// fmt.Printf("Block ID: %s\n", block.ID)
		for _, rt := range block.Paragraph.RichText {
			// This example directly writes the content. You might want to format it as valid Markdown.
			_, err := file.WriteString(rt.Text.Content + "\n")
			// When we return a block that is empty then we'll Add two newlines
			if err != nil {
				return fmt.Errorf("error writing to markdown file: %w", err)
			}
			fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
		}
	}

	return nil // Return nil if everything was successful
}
